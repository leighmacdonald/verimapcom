package web

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/leighmacdonald/verimapcom/store"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func successFlash(c *gin.Context, msg string, path string) {
	flash(c, lSuccess, msg)
	c.Redirect(http.StatusFound, path)
}

func abortFlash(c *gin.Context, msg string, path string) {
	flash(c, lError, msg)
	c.Redirect(http.StatusFound, path)
}

func abortFlashErr(c *gin.Context, msg string, path string, err error) {
	abortFlash(c, msg, path)
	log.Error(err)
}

func (w *Web) route(name pageName, methods ...string) string {
	method := "GET"
	if len(methods) > 0 {
		method = methods[0]
	}
	r, found := w.pages[method][name]
	if !found {
		return "/"
	}
	return r.Path
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(b), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func adminMiddleWare(w *Web) gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := c.Get("person")
		if !ok || p.(store.Person).AgencyID != 1 {
			abortFlash(c, "You must login to access this resource", w.route(login))
			return
		}
		c.Next()
	}
}

func sessionMiddleWare(ctx context.Context, db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		guest := store.Person{
			PersonID:  0,
			FirstName: "Guest",
		}
		var p store.Person
		v := s.Get("person_id")
		if v != nil {
			pId, ok := v.(int32)
			if ok {
				if err := store.LoadPersonByID(ctx, db, pId, &p); err != nil {
					log.Errorf("Failed to load persons session user: %v", err)
					p = guest
				}
			} else {
				// Delete the bad value
				s.Delete("person_id")
				if err := s.Save(); err != nil {
					log.Errorf("Failed to save session")
				}
			}
		} else {
			p = guest
		}
		c.Set("person", p)
		c.Next()
	}
}

type pageName string

const (
	about              pageName = "about"
	adminAgencies      pageName = "admin_agencies"
	adminPeople        pageName = "admin_people"
	adminPeopleDelete  pageName = "admin_people_delete"
	adminPeopleEdit    pageName = "admin_people_edit"
	background         pageName = "background"
	connect            pageName = "connect"
	connectSendMessage pageName = "connect_send_message"
	downloads          pageName = "downloads"
	downloadFile       pageName = "downloads_file"
	emergency          pageName = "emergency"
	environmental      pageName = "environmental"
	err                pageName = "error"
	example            pageName = "example"
	examples           pageName = "examples"
	firetracker        pageName = "firetracker"
	home               pageName = "home"
	infrastructure     pageName = "infrastructure"
	login              pageName = "login"
	logout             pageName = "logout"
	mission            pageName = "mission"
	missionEvents      pageName = "mission_events"
	missions           pageName = "missions"
	missionsCreate     pageName = "missions_create"
	partners           pageName = "partners"
	profile            pageName = "profile"
	profileCreate      pageName = "profile_create"
	services           pageName = "services"
	technology         pageName = "technology"
	upload             pageName = "upload"
	uploads            pageName = "uploads"
	wildfire           pageName = "wildfire"
)

var ErrInvalidUser = errors.New("Invalid user")

func (w *Web) currentPerson(c *gin.Context) (store.Person, error) {
	p, found := c.Get("person")
	if !found {
		return store.Person{FirstName: "Guest"}, ErrInvalidUser
	}
	person, ok := p.(store.Person)
	if !ok {
		log.Warnf("Count not cast store.Person from session")
		return store.Person{FirstName: "Guest"}, ErrInvalidUser
	}
	return person, nil
}

func (w *Web) defaultM(c *gin.Context, page pageName) M {
	p, _ := w.currentPerson(c)
	m := M{
		"person": p,
	}
	headers, found := w.headers[page]
	if found {
		m["header"] = headers
	}
	s := sessions.Default(c)
	flashes := s.Flashes()
	if len(flashes) > 0 {
		m["alerts"] = flashes
		pid := s.Get("person_id")
		s.Clear()
		if err := s.Save(); err != nil {
			log.Errorf("Failed to save clear session call: %v", err)
			return m
		}
		s.Set("person_id", pid)
		if err := s.Save(); err != nil {
			log.Errorf("Failed to save person_id to updated session: %v", err)
			return m
		}
	}
	return m
}

func (w *Web) newTmpl(files ...string) *template.Template {
	var tFuncMap = template.FuncMap{
		"img": img,
		"md":  md,
		"route": func(p pageName) template.HTML {
			return template.HTML(w.route(p))
		},
		"icon": func(class string) template.HTML {
			return template.HTML(fmt.Sprintf(`<i class="%s"></i>`, class))
		},
		"currentYear": func() template.HTML {
			return template.HTML(fmt.Sprintf("%d", time.Now().UTC().Year()))
		},
		"human": func(size int64) template.HTML {
			return template.HTML(ByteCountSI(size))
		},
		"datetime": func(t time.Time) template.HTML {
			return template.HTML(t.Format(time.RFC822))
		},
		"fmtFloat": func(f float64, size int) template.HTML {
			ft := fmt.Sprintf("%%.%df", size)
			return template.HTML(fmt.Sprintf(ft, f))
		},
	}
	tmpl, err := template.New("layout").Funcs(tFuncMap).ParseFiles(files...)
	if err != nil {
		log.Panicf("Failed to load template: %v", err)
	}
	return tmpl
}

func (w *Web) render(c *gin.Context, t pageName, a M) {
	var buf bytes.Buffer
	tmpl := w.tmpl[t]
	if err := tmpl.ExecuteTemplate(&buf, "layout", a); err != nil {
		log.Errorf("Failed to execute template: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Data(200, gin.MIMEHTML, buf.Bytes())
}

// HTTPOpts is used to configure a http.Server instance
type HTTPOpts struct {
	ListenAddr     string
	UseTLS         bool
	Handler        http.Handler
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	TLSConfig      *tls.Config
}

// DefaultHTTPOpts returns a default set of options for http.Server instances
func DefaultHTTPOpts() *HTTPOpts {
	addr := viper.GetString("listen_http")
	return &HTTPOpts{
		ListenAddr:     addr,
		UseTLS:         false,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      nil,
	}
}

// NewHTTPServer will configure and return a *http.Server suitable for serving requests.
// This should be used over the default ListenAndServe options as they do not set certain
// parameters, notably timeouts, which can negatively effect performance.
func NewHTTPServer(opts *HTTPOpts) *http.Server {
	var tlsCfg *tls.Config
	if opts.UseTLS && opts.TLSConfig == nil {
		tlsCfg = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
	} else {
		tlsCfg = nil
	}
	srv := &http.Server{
		Addr:           opts.ListenAddr,
		Handler:        opts.Handler,
		TLSConfig:      tlsCfg,
		ReadTimeout:    opts.ReadTimeout,
		WriteTimeout:   opts.WriteTimeout,
		MaxHeaderBytes: opts.MaxHeaderBytes,
	}
	return srv
}

type Web struct {
	Handler    http.Handler
	db         *pgxpool.Pool
	ctx        context.Context
	client     *http.Client
	uploadPath string
	tmpl       map[pageName]*template.Template
	pages      map[string]map[pageName]*Page
	headers    map[pageName]*Header
	wsClientMu *sync.RWMutex
}

func (w *Web) Setup() error {
	if !store.Exists(w.uploadPath) {
		if err := os.MkdirAll(w.uploadPath, 0755); err != nil {
			log.Fatalf("Failed to create upload directory: %v", err)
		}
	}
	var a0 store.Agency
	var p0 store.Person
	if err := store.LoadAgency(w.ctx, w.db, 1, &a0); err != nil {
		if err.Error() != pgx.ErrNoRows.Error() {
			log.Fatalf("Unhandled db error loadAgency: %v", err)
		}
		a0.CreatedOn = time.Now()
		a0.AgencyName = "Verimap Plus Inc."
		if err := store.SaveAgency(w.ctx, w.db, &a0); err != nil {
			log.Fatalf("Failed to load initial agency: %v", err)
		}
	}
	if err := store.LoadPersonByID(w.ctx, w.db, 1, &p0); err != nil {
		if err.Error() != pgx.ErrNoRows.Error() {
			log.Fatalf("Unhandled db error loadPersonByID: %v", err)
		}
		p0.CreatedOn = time.Now()
		p0.FirstName = "Leigh"
		p0.LastName = "MacDonald"
		p0.AgencyID = a0.AgencyID
		pw, err := HashPassword("temp")
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}
		p0.PasswordHash = pw
		p0.Email = "leigh.macdonald@gmail.com"
		if err := store.SavePerson(w.ctx, w.db, &p0); err != nil {
			log.Fatalf("Failed to create initial user: %v", err)
		}
	}
	return nil
}

func (w *Web) Close() {
	w.db.Close()
}

func New(ctx context.Context, redisHost string) *Web {
	s, err := redis.NewStoreWithDB(10, "tcp",
		redisHost, "", "5", []byte("temp"))
	if err != nil {
		log.Fatalf("Could not connect to redis: %v", err)
	}

	var templateFiles []string
	root := "templates"
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".gohtml") {
			if !strings.Contains(path, "/layouts/") && !strings.Contains(path, "/partials") {
				templateFiles = append(templateFiles, info.Name())
			}
		}
		return nil
	}); err != nil {
		log.Fatalf("Failed to read templates: %v", err)
	}

	r := gin.Default()
	w := &Web{
		r,
		store.MustConnectDB(ctx),
		ctx,
		&http.Client{
			Timeout: time.Second * 20,
		},
		"./uploads",
		make(map[pageName]*template.Template),
		make(map[string]map[pageName]*Page),
		make(map[pageName]*Header),
		&sync.RWMutex{},
	}
	w.setup()

	staticPath := viper.GetString("static_path")
	if staticPath == "" {
		staticPath = "frontend/dist"
	}
	r.Static("/dist", staticPath)
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	var newPagesSet = func(path string) []string {
		return []string{
			fmt.Sprintf("templates/%s.gohtml", path),
			"templates/partials/page_header.gohtml",
			"templates/layouts/layout.gohtml",
		}
	}
	sesh := r.Group("", sessions.Sessions("vmsesh", s), sessionMiddleWare(w.ctx, w.db))
	admin := sesh.Group("", adminMiddleWare(w))

	sesh.GET("/ws", func(c *gin.Context) {
		wsHandler(w, c)
	})

	for _, p := range templateFiles {
		pageN := strings.ReplaceAll(p, ".gohtml", "")
		w.tmpl[pageName(pageN)] = w.newTmpl(newPagesSet(pageN)...)
	}
	for method, handlers := range w.pages {
		for page, p := range handlers {
			if p.Handler == nil {
				p.Handler = w.simple(page)
			}
			fn := sesh
			if p.Admin {
				fn = admin
			}
			switch method {
			case "GET":
				fn.GET(p.Path, p.Handler)
			case "POST":
				fn.POST(p.Path, p.Handler)
			}
		}
	}
	return w
}
func (w *Web) page(name pageName) *Page {
	return w.pages["GET"][name]
}

func (w *Web) setup() {
	pages := map[string]map[pageName]*Page{
		"GET": {
			adminAgencies:   {Path: "/admin/agencies", Admin: true, Handler: w.getAdminAgencies},
			adminPeople:     {Path: "/admin/people", Admin: true, Handler: w.getAdminPeople},
			adminPeopleEdit: {Path: "/admin/people/:person_id", Admin: true, Handler: w.getAdminPeopleEdit},
			about:           {Path: "/about"},
			background:      {Path: "/about/background", Handler: w.getBackground},
			connect:         {Path: "/connect"},
			downloads:       {Path: "/downloads", Handler: w.getDownloads},
			downloadFile:    {Path: "/download/:file_id", Handler: w.getFile},
			emergency:       {Path: "/services/emergency"},
			environmental:   {Path: "/services/environmental"},
			err:             {Path: "/error"},
			example:         {Path: "/example/:example_id", Handler: w.getExample},
			examples:        {Path: "/examples", Handler: w.getExamples},
			firetracker:     {Path: "/firetracker", Handler: w.getFireTracker},
			home:            {Path: "/", Handler: w.getHome},
			infrastructure:  {Path: "/services/infrastructure"},
			login:           {Path: "/login", Handler: w.getLogin},
			logout:          {Path: "/logout", Handler: w.getLogout},
			mission:         {Path: "/mission/:mission_id", Handler: w.getMission},
			missionEvents:   {Path: "/mission/:mission_id/events", Handler: w.getMissionEvents},
			missions:        {Path: "/missions", Handler: w.getMissions},
			missionsCreate:  {Path: "/missions/create", Handler: w.getMissionsCreate},
			partners:        {Path: "/about/partners", Handler: w.getPartners},
			profile:         {Path: "/profile", Handler: w.getProfile},
			profileCreate:   {Path: "/profile/create", Handler: w.getProfileCreate},
			services:        {Path: "/services"},
			technology:      {Path: "/innovation/technology"},
			upload:          {Path: "/upload", Handler: w.getUpload},
			uploads:         {Path: "/uploads", Handler: w.getUploads},
			wildfire:        {Path: "/services/wildfire"},
		},
		"POST": {
			upload:             {Path: "/upload", Handler: w.postUpload},
			login:              {Path: "/login", Handler: w.postLogin},
			profile:            {Path: "/profile", Handler: w.postProfile},
			missionsCreate:     {Path: "/mission/create", Handler: w.postMission},
			connectSendMessage: {Path: "/connect/send", Handler: w.sendConnectMessage},
			adminPeopleEdit:    {Path: "/admin/people/:person_id", Handler: w.postAdminPeopleEdit, Admin: true},
			adminPeopleDelete:  {Path: "/admin/people_delete/:person_id", Handler: w.postAdminPeopleDelete, Admin: true},
			adminPeople:        {Path: "/admin/people", Handler: w.postAdminPeopleCreate, Admin: true},
			adminAgencies:      {Path: "/admin/agencies", Handler: w.postAdminAgenciesCreate, Admin: true},
		},
	}

	w.pages = pages

	headers := map[pageName]*Header{
		services: {
			Img:         "/dist/images/golden_gate_shore.png",
			Name:        "Services",
			Breadcrumbs: []*Page{w.page(home), w.page(services)},
		},
		wildfire: {
			Img:         "/dist/images/fire_fighters_12_2.png",
			Name:        "Wildfire Mapping Services",
			Breadcrumbs: []*Page{w.page(home), w.page(services), w.page(wildfire)},
		},
		emergency: {
			Img:         "/dist/images/header_emergency.png",
			Name:        "Emergency Response Management",
			Breadcrumbs: []*Page{w.page(home), w.page(services), w.page(emergency)},
		},
		environmental: {
			Img:         "/dist/images/false_colour_dem.png",
			Name:        "Environmental",
			Breadcrumbs: []*Page{w.page(home), w.page(services), w.page(environmental)},
		},
		infrastructure: {
			Img:         "/dist/images/barrels.png",
			Name:        "Infrastructure",
			Breadcrumbs: []*Page{w.page(home), w.page(services), w.page(infrastructure)},
		},

		technology: {
			Img:         "/dist/images/contours.png",
			Name:        "Technology",
			Breadcrumbs: []*Page{w.page(home), w.page(technology)},
		},
		examples: {
			Img:         "/dist/images/header_solar.png",
			Name:        "Example Datasets",
			Breadcrumbs: []*Page{w.page(home), w.page(examples)},
		},
		example: {
			Img:         "/dist/images/header_solar.png",
			Name:        "Example Dataset",
			Breadcrumbs: []*Page{w.page(home), w.page(examples), w.page(example)},
		},
		firetracker: {
			Img:         "/dist/images/fire_fighters_12_2.png",
			Name:        "Global Fire Tracker",
			Breadcrumbs: []*Page{w.page(home), w.page(firetracker)},
		},
		background: {
			Img:         "/dist/images/header_emergency.png",
			Name:        "Background",
			Breadcrumbs: []*Page{w.page(home), w.page(about), w.page(background)},
		},
		partners: {
			Img:         "/dist/images/golden_gate_shore.png",
			Name:        "Partners",
			Breadcrumbs: []*Page{w.page(home), w.page(about), w.page(partners)},
		},
		connect: {
			Img:         "/dist/images/header_contact.png",
			Name:        "Connect With Us",
			Breadcrumbs: []*Page{w.page(home), w.page(connect)},
		},
		login: {
			Img:         "/dist/images/false_colour_dem.png",
			Name:        "User Login",
			Breadcrumbs: []*Page{w.page(home), w.page(login)},
		},
		logout: {
			Img:         "/dist/images/false_colour_dem.png",
			Name:        "User Logout",
			Breadcrumbs: []*Page{w.page(home), w.page(logout)},
		},
		profile: {
			Img:         "/dist/images/header_contact.png",
			Name:        "Person Profile",
			Breadcrumbs: []*Page{w.page(home), w.page(profile)},
		},
	}
	w.headers = headers

}
func init() {
	gob.Register(Flash{})
}
