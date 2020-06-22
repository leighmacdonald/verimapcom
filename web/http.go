package web

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/leighmacdonald/verimapcom/web/store"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type exampleKML struct {
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	Hash             string      `json:"hash"`
	Sha256           string      `json:"sha256"`
	Ext              string      `json:"ext"`
	Mime             string      `json:"mime"`
	Size             string      `json:"size"`
	URL              string      `json:"url"`
	Provider         string      `json:"provider"`
	ProviderMetadata interface{} `json:"provider_metadata"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

type examplePage struct {
	ID           int               `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Latitude     float64           `json:"latitude"`
	Longitude    float64           `json:"longitude"`
	Zoom         int               `json:"zoom"`
	Public       bool              `json:"public"`
	Stats        string            `json:"stats"`
	StatsMap     map[string]string `json:"-"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Layer        string            `json:"layer"`
	ZoomMin      int               `json:"zoom_min"`
	ZoomMax      int               `json:"zoom_max"`
	VectorLayers interface{}       `json:"vector_layers"`
	Kml          exampleKML        `json:"kml,omitempty"`
}

type Showcase struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	LinkText  string    `json:"link_text"`
	LinkURL   string    `json:"link_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Order     int       `json:"order"`
	Image     struct {
		ID               int         `json:"id"`
		Name             string      `json:"name"`
		Hash             string      `json:"hash"`
		Sha256           string      `json:"sha256"`
		Ext              string      `json:"ext"`
		Mime             string      `json:"mime"`
		Size             string      `json:"size"`
		URL              string      `json:"url"`
		Provider         string      `json:"provider"`
		ProviderMetadata interface{} `json:"provider_metadata"`
		CreatedAt        time.Time   `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
	} `json:"image"`
}

type FireWatch struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Published time.Time `json:"published"`
	Body      string    `json:"body"`
	User      struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Provider  string    `json:"provider"`
		Confirmed bool      `json:"confirmed"`
		Blocked   bool      `json:"blocked"`
		Role      int       `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	LinkText  string      `json:"link_text"`
	LinkURL   interface{} `json:"link_url"`
	Thumbnail struct {
		ID               int         `json:"id"`
		Name             string      `json:"name"`
		Hash             string      `json:"hash"`
		Sha256           string      `json:"sha256"`
		Ext              string      `json:"ext"`
		Mime             string      `json:"mime"`
		Size             string      `json:"size"`
		URL              string      `json:"url"`
		Provider         string      `json:"provider"`
		ProviderMetadata interface{} `json:"provider_metadata"`
		CreatedAt        time.Time   `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
	} `json:"thumbnail"`
	Gallery []interface{} `json:"gallery"`
}

func apiGetFireWatches(ctx context.Context, client *http.Client, count int) ([]FireWatch, error) {
	var resp []FireWatch
	u := fmt.Sprintf("https://cms.verimap.com/firewatches?_sort=published:desc&_limit=%d", count)
	if err := get(ctx, client, u, &resp); err != nil {
		return nil, errors.Wrapf(err, "Failed to make get request")
	}
	return resp, nil
}

func apiGetShowcases(ctx context.Context, client *http.Client) ([]Showcase, error) {
	var resp []Showcase
	if err := get(ctx, client, "https://cms.verimap.com/showcases?_sort=order", &resp); err != nil {
		return nil, errors.Wrapf(err, "Failed to make get request")
	}
	return resp, nil
}

func apiGetExample(ctx context.Context, client *http.Client, ID int) (examplePage, error) {
	var resp []examplePage
	url := fmt.Sprintf("https://cms.verimap.com/examples?public=true&id=%d", ID)
	if err := get(ctx, client, url, &resp); err != nil {
		return examplePage{}, errors.Wrapf(err, "Failed to make get request")
	}
	for _, page := range resp {
		m := make(map[string]string)
		for _, row := range strings.Split(page.Stats, "\n") {
			cols := strings.SplitN(row, "|", 2)
			if len(cols) == 2 {
				m[cols[0]] = cols[1]
			}
		}
		page.StatsMap = m
		return page, nil
	}
	return examplePage{}, errors.New("Unknown result")
}

func apiGetExamples(ctx context.Context, client *http.Client) ([]examplePage, error) {
	var resp []examplePage
	if err := get(ctx, client, "https://cms.verimap.com/examples?public=true", &resp); err != nil {
		return nil, errors.Wrapf(err, "Failed to make get request")
	}
	for i, page := range resp {
		m := make(map[string]string)
		for _, row := range strings.Split(page.Stats, "\n") {
			cols := strings.SplitN(row, "|", 2)
			if len(cols) == 2 {
				m[cols[0]] = cols[1]
			}
		}
		resp[i].StatsMap = m
	}
	return resp, nil
}

type Level string

const (
	lSuccess Level = "success"
	lWarning Level = "warning"
	lError   Level = "alert"
)

type Flash struct {
	Level   Level  `json:"level"`
	Message string `json:"message"`
}

func formFloatDefault(c *gin.Context, field string, def float64) float64 {
	f, err := strconv.ParseFloat(c.PostForm(field), 64)
	if err != nil {
		return def
	}
	return f
}

func flash(c *gin.Context, level Level, msg string) {
	s := sessions.Default(c)
	s.AddFlash(Flash{
		Level:   level,
		Message: msg,
	})
	if err := s.Save(); err != nil {
		log.Errorf("failed to save flash")
	}
}

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

func (w *Web) route(name pageName) string {
	r, found := w.pages[name]
	if !found {
		return "/"
	}
	return r.Path
}

func get(ctx context.Context, client *http.Client, url string, recv interface{}) error {
	c, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	req, err := http.NewRequestWithContext(c, "GET", url, nil)
	if err != nil {
		return errors.Wrapf(err, "Failed to create request")
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Failed to perform request")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "Failed to read response body")
	}
	if err := json.Unmarshal(b, recv); err != nil {
		return errors.Wrapf(err, "Failed to decode json response")
	}
	return nil
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
			pId, ok := v.(int)
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
	about             pageName = "about"
	adminAgencies     pageName = "admin_agencies"
	adminMissions     pageName = "admin_missions"
	adminPeople       pageName = "admin_people"
	adminPeopleDelete pageName = "admin_people_delete"
	adminPeopleEdit   pageName = "admin_people_edit"
	background        pageName = "background"
	connect           pageName = "connect"
	download          pageName = "download"
	downloads         pageName = "downloads"
	emergency         pageName = "emergency"
	environmental     pageName = "environmental"
	err               pageName = "error"
	example           pageName = "example"
	examples          pageName = "examples"
	firetracker       pageName = "firetracker"
	home              pageName = "home"
	infrastructure    pageName = "infrastructure"
	login             pageName = "login"
	logout            pageName = "logout"
	mission           pageName = "mission"
	missions          pageName = "missions"
	missionsCreate    pageName = "missions_create"
	partners          pageName = "partners"
	profile           pageName = "profile"
	services          pageName = "services"
	technology        pageName = "technology"
	upload            pageName = "upload"
	uploads           pageName = "uploads"
	wildfire          pageName = "wildfire"
)

type Header struct {
	Img         string
	Name        string
	Breadcrumbs []*Page
}

type Page struct {
	Name    string
	Path    string
	Handler gin.HandlerFunc
	Admin   bool
}

type M map[string]interface{}

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

func img(src, alt string, trans bool) template.HTML {
	var s strings.Builder
	if trans {
		s.WriteString(`<div class="callout_trans">`)
	} else {
		s.WriteString(`<div class="callout">`)
	}
	s.WriteString(fmt.Sprintf(`
	<figure>
	    <img src="%s" alt="%s">
	</figure>
	</div>`, src, alt))
	return template.HTML(s.String())
}

func (w *Web) urlFor(page string) template.HTML {
	for _, r := range w.pages {
		if page == r.Name {
			return template.HTML(r.Path)
		}
	}
	log.Panicf("Unknown Page path: %s", page)
	return "#"
}

func md(data string) template.HTML {
	out := markdown.ToHTML([]byte(data), nil, nil)
	return template.HTML(out)
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
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
	addr, found := os.LookupEnv("LISTEN")
	if !found {
		addr = ":8080"
	}
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
	pages      map[pageName]*Page
	headers    map[pageName]*Header
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

func New(ctx context.Context) *Web {
	redisHost, found := os.LookupEnv("REDIS_HOST")
	if !found {
		redisHost = "localhost:6379"
	}
	s, err := redis.NewStoreWithDB(10, "tcp",
		redisHost, "", "5", []byte("temp"))
	if err != nil {
		log.Fatalf("Could not connect to redis: %v", err)
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
		make(map[pageName]*Page),
		make(map[pageName]*Header),
	}
	w.makePagesHeaders()

	r.Static("/dist", "dist")
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
	admin.POST(w.route(adminPeopleEdit), w.postAdminPeopleEdit)
	admin.POST("/admin/people_delete/:person_id", w.postAdminPeopleDelete)
	admin.POST(w.route(adminPeople), w.postAdminPeopleCreate)
	admin.POST(w.route(adminAgencies), w.postAdminAgenciesCreate)
	for page, p := range w.pages {
		w.tmpl[page] = w.newTmpl(newPagesSet(p.Name)...)
		if p.Handler == nil {
			p.Handler = w.simple(page)
		}
		if p.Admin {
			admin.GET(p.Path, p.Handler)
		} else {
			sesh.GET(p.Path, p.Handler)
		}
	}
	sesh.GET("/download/:file_id", w.getFile)
	sesh.POST(w.route(upload), w.postUpload)
	sesh.POST(w.route(login), w.postLogin)
	sesh.POST(w.route(profile), w.postProfile)
	sesh.POST(w.route(missionsCreate), w.postMission)

	return w
}

func (w *Web) makePagesHeaders() {
	pages := map[pageName]*Page{
		adminAgencies:   {Name: "admin_agencies", Path: "/admin/agencies", Admin: true, Handler: w.getAdminAgencies},
		adminMissions:   {Name: "admin_missions", Path: "/admin/missions", Admin: true, Handler: w.getAdminMissions},
		adminPeople:     {Name: "admin_people", Path: "/admin/people", Admin: true, Handler: w.getAdminPeople},
		adminPeopleEdit: {Name: "admin_people_edit", Path: "/admin/people/:person_id", Admin: true, Handler: w.getAdminPeopleEdit},
		about:           {Name: "about", Path: "/about"},
		background:      {Name: "background", Path: "/about/background"},
		connect:         {Name: "connect", Path: "/connect"},
		downloads:       {Name: "downloads", Path: "/downloads", Handler: w.getDownloads},
		emergency:       {Name: "emergency", Path: "/services/emergency"},
		environmental:   {Name: "environmental", Path: "/services/environmental"},
		err:             {Name: "error", Path: "/error"},
		example:         {Name: "example", Path: "/example/:example_id", Handler: w.getExample},
		examples:        {Name: "examples", Path: "/examples", Handler: w.getExamples},
		firetracker:     {Name: "firetracker", Path: "/firetracker", Handler: w.getFireTracker},
		home:            {Name: "home", Path: "/", Handler: w.getHome},
		infrastructure:  {Name: "infrastructure", Path: "/services/infrastructure"},
		login:           {Name: "login", Path: "/login", Handler: w.getLogin},
		logout:          {Name: "logout", Path: "/logout", Handler: w.getLogout},
		mission:         {Name: "mission", Path: "/mission/:mission_id", Handler: w.getMission},
		missions:        {Name: "missions", Path: "/missions", Handler: w.getMissions},
		missionsCreate:  {Name: "missions_create", Path: "/missions/create", Handler: w.getMissionsCreate},
		partners:        {Name: "partners", Path: "/about/partners", Handler: w.getPartners},
		profile:         {Name: "profile", Path: "/profile", Handler: w.getProfile},
		services:        {Name: "services", Path: "/services"},
		technology:      {Name: "technology", Path: "/innovation/technology"},
		upload:          {Name: "upload", Path: "/upload", Handler: w.getUpload},
		uploads:         {Name: "uploads", Path: "/uploads", Handler: w.getUploads},
		wildfire:        {Name: "wildfire", Path: "/services/wildfire"},
	}
	w.pages = pages

	headers := map[pageName]*Header{
		services: {
			Img:         "/dist/assets/golden_gate_shore.png",
			Name:        "Services",
			Breadcrumbs: []*Page{pages[home], pages[services]},
		},
		wildfire: {
			Img:         "/dist/assets/fire_fighters_12_2.png",
			Name:        "Wildfire Mapping Services",
			Breadcrumbs: []*Page{pages[home], pages[services], pages[wildfire]},
		},
		emergency: {
			Img:         "/dist/assets/header_emergency.png",
			Name:        "Emergency Response Management",
			Breadcrumbs: []*Page{pages[home], pages[services], pages[emergency]},
		},
		environmental: {
			Img:         "/dist/assets/false_colour_dem.png",
			Name:        "Environmental",
			Breadcrumbs: []*Page{pages[home], pages[services], pages[environmental]},
		},
		infrastructure: {
			Img:         "/dist/assets/barrels.png",
			Name:        "Infrastructure",
			Breadcrumbs: []*Page{pages[home], pages[services], pages[infrastructure]},
		},

		technology: {
			Img:         "/dist/assets/contours.png",
			Name:        "Technology",
			Breadcrumbs: []*Page{pages[home], pages[technology]},
		},
		examples: {
			Img:         "/dist/assets/header_solar.png",
			Name:        "Example Datasets",
			Breadcrumbs: []*Page{pages[home], pages[examples]},
		},
		example: {
			Img:         "/dist/assets/header_solar.png",
			Name:        "Example Dataset",
			Breadcrumbs: []*Page{pages[home], pages[examples], pages[example]},
		},
		firetracker: {
			Img:         "/dist/assets/fire_fighters_12_2.png",
			Name:        "Global Fire Tracker",
			Breadcrumbs: []*Page{pages[home], pages[firetracker]},
		},
		background: {
			Img:         "/dist/assets/header_emergency.png",
			Name:        "Background",
			Breadcrumbs: []*Page{pages[home], pages[about], pages[background]},
		},
		partners: {
			Img:         "/dist/assets/golden_gate_shore.png",
			Name:        "Partners",
			Breadcrumbs: []*Page{pages[home], pages[about], pages[partners]},
		},
		connect: {
			Img:         "/dist/assets/header_contact.png",
			Name:        "Connect With Us",
			Breadcrumbs: []*Page{pages[home], pages[connect]},
		},
		login: {
			Img:         "/dist/assets/false_colour_dem.png",
			Name:        "User Login",
			Breadcrumbs: []*Page{pages[home], pages[login]},
		},
		logout: {
			Img:         "/dist/assets/false_colour_dem.png",
			Name:        "User Logout",
			Breadcrumbs: []*Page{pages[home], pages[logout]},
		},
		profile: {
			Img:         "/dist/assets/header_contact.png",
			Name:        "Person Profile",
			Breadcrumbs: []*Page{pages[home], pages[profile]},
		},
	}
	w.headers = headers

}
func init() {
	gob.Register(Flash{})
}
