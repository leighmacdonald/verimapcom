package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
)

type pageName string

const (
	about          pageName = "about"
	adminAgencies  pageName = "admin_agencies"
	adminMissions  pageName = "admin_missions"
	adminPeople    pageName = "admin_people"
	background     pageName = "background"
	connect        pageName = "connect"
	emergency      pageName = "emergency"
	environmental  pageName = "environmental"
	err            pageName = "error"
	example        pageName = "example"
	examples       pageName = "examples"
	firetracker    pageName = "firetracker"
	home           pageName = "home"
	infrastructure pageName = "infrastructure"
	login          pageName = "login"
	logout         pageName = "logout"
	missions       pageName = "missions"
	partners       pageName = "partners"
	services       pageName = "services"
	technology     pageName = "technology"
	wildfire       pageName = "wildfire"
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

var (
	tmpl    map[pageName]*template.Template
	pages   map[pageName]*Page
	headers map[pageName]*Header
	client  *http.Client
	ctx     context.Context
	dbpool  *pgxpool.Pool
)

type M map[string]interface{}

func defaultM(c *gin.Context, page pageName) M {
	p, found := c.Get("person")
	if !found {
		p = Person{FirstName: "Guest"}
	}
	m := M{
		"person": p.(Person),
	}
	headers, found := headers[page]
	if found {
		m["header"] = headers
	}
	s := sessions.Default(c)
	flashes := s.Flashes()
	if len(flashes) > 0 {
		m["flashes"] = flashes
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

func urlFor(page string) template.HTML {
	for _, r := range pages {
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

var tFuncMap = template.FuncMap{
	"img":     img,
	"url_for": urlFor,
	"md":      md,
}

func newTmpl(files ...string) *template.Template {
	tmpl, err := template.New("layout").Funcs(tFuncMap).ParseFiles(files...)
	if err != nil {
		log.Panicf("Failed to load template: %v", err)
	}
	return tmpl
}

func render(c *gin.Context, t pageName, a M) {
	var buf bytes.Buffer
	tmpl := tmpl[t]
	if err := tmpl.ExecuteTemplate(&buf, "layout", a); err != nil {
		log.Errorf("Failed to execute template: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Data(200, gin.MIMEHTML, buf.Bytes())
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(b), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func adminMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		p, ok := c.Get("person")
		if !ok || p.(Person).AgencyID != 1 {
			s := sessions.Default(c)
			s.AddFlash("You must login to access this resource")
			if err := s.Save(); err != nil {
				log.Errorf("Failed to save flash in admin mw: %v", err)
			}
			c.Redirect(http.StatusFound, pages[login].Path)
			return
		}
		c.Next()
	}
}

func sessionMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		guest := Person{
			PersonID:  0,
			FirstName: "Guest",
		}
		var p Person
		v := s.Get("person_id")
		if v != nil {
			pId, ok := v.(int)
			if ok {
				if err := loadPersonByID(pId, &p); err != nil {
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

func main() {
	ctx := context.Background()
	dbpool = mustConnectDB(ctx)
	defer dbpool.Close()

	var a0 Agency
	var p0 Person
	if err := loadAgency(1, &a0); err != nil {
		if err.Error() != pgx.ErrNoRows.Error() {
			log.Fatalf("Unhandled db error loadAgency: %v", err)
		}
		a0.CreatedOn = time.Now()
		a0.AgencyName = "Verimap Plus Inc."
		if err := saveAgency(&a0); err != nil {
			log.Fatalf("Failed to load initial agency: %v", err)
		}
	}
	if err := loadPersonByID(1, &p0); err != nil {
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
		if err := savePerson(&p0); err != nil {
			log.Fatalf("Failed to create initial user: %v", err)
		}
	}
	redisHost, found := os.LookupEnv("REDIS_HOST")
	if !found {
		redisHost = "172.16.1.4:6379"
	}
	store, err := redis.NewStoreWithDB(10, "tcp",
		redisHost, "", "5", []byte("temp"))
	if err != nil {
		log.Fatalf("Could not connect to redis: %v", err)
	}
	r := gin.Default()

	r.Static("/dist", "dist")
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	var newPagesSet = func(path string) []string {
		return []string{
			fmt.Sprintf("templates/%s.gohtml", path),
			"templates/partials/page_header.gohtml",
			"templates/layouts/layout.gohtml",
		}
	}
	sesh := r.Group("", sessions.Sessions("vmsesh", store), sessionMiddleWare())
	admin := sesh.Group("", adminMiddleWare())
	for page, p := range pages {
		tmpl[page] = newTmpl(newPagesSet(p.Name)...)
		if p.Handler == nil {
			p.Handler = simple(page)
		}
		if p.Admin {
			admin.GET(p.Path, p.Handler)
		} else {
			sesh.GET(p.Path, p.Handler)
		}
	}
	sesh.POST(pages[login].Path, postLogin)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err := r.Run(); err != nil {
		log.Errorf("Shutdown unclean: %v", err)
	}
}

func init() {
	ctx = context.Background()
	client = &http.Client{
		Timeout: time.Second * 20,
	}
	tmpl = make(map[pageName]*template.Template)
	pages = map[pageName]*Page{
		adminAgencies:  {Name: "admin_agencies", Path: "/admin/agencies", Admin: true, Handler: getAdminAgencies},
		adminMissions:  {Name: "admin_missions", Path: "/admin/missions", Admin: true, Handler: getAdminMissions},
		adminPeople:    {Name: "admin_people", Path: "/admin/people", Admin: true, Handler: getAdminPeople},
		about:          {Name: "about", Path: "/about"},
		background:     {Name: "background", Path: "/about/background"},
		connect:        {Name: "connect", Path: "/connect"},
		emergency:      {Name: "emergency", Path: "/services/emergency"},
		environmental:  {Name: "environmental", Path: "/services/environmental"},
		err:            {Name: "error", Path: "/error"},
		example:        {Name: "example", Path: "/example/:example_id", Handler: getExample},
		examples:       {Name: "examples", Path: "/examples", Handler: getExamples},
		firetracker:    {Name: "firetracker", Path: "/firetracker", Handler: getFireTracker},
		home:           {Name: "home", Path: "/", Handler: getHome},
		infrastructure: {Name: "infrastructure", Path: "/services/infrastructure"},
		login:          {Name: "login", Path: "/login", Handler: getLogin},
		logout:         {Name: "logout", Path: "/logout", Handler: getLogout},
		missions:       {Name: "missions", Path: "/missions", Handler: getMissions},
		partners:       {Name: "partners", Path: "/about/partners", Handler: getPartners},
		services:       {Name: "services", Path: "/services"},
		technology:     {Name: "technology", Path: "/innovation/technology"},
		wildfire:       {Name: "wildfire", Path: "/services/wildfire"},
	}
	headers = map[pageName]*Header{
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
	}
}
