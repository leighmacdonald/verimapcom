package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/leighmacdonald/verimapcom/web/store"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (w *Web) simple(page pageName) gin.HandlerFunc {
	return func(c *gin.Context) {
		w.render(c, page, w.defaultM(c, page))
	}
}

func (w *Web) getLogout(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("person_id")
	if v == nil {
		c.Redirect(http.StatusFound, w.route(login))
		return
	}
	_, ok := v.(int)
	if !ok {
		c.Redirect(http.StatusFound, w.route(login))
		return
	}
	session.Clear()
	if err := session.Save(); err != nil {
		log.Errorf("failed to clear user session on logout: %v", err)
	}
	c.Redirect(http.StatusFound, "/")
}

func (w *Web) postLogin(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("person_id")
	if v != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	email := c.PostForm("email")
	password := c.PostForm("password")
	var p store.Person
	if err := store.LoadPersonByEmail(w.ctx, w.db, email, &p); err != nil {
		abortFlashErr(c, "Invalid username or password", w.route(login), err)
		return
	}
	if !CheckPasswordHash(password, p.PasswordHash) {
		abortFlash(c, "Invalid username or password", w.route(login))
		return
	}
	session.Set("person_id", p.PersonID)
	session.AddFlash(Flash{
		Level:   lSuccess,
		Message: "Logged in successfully",
	})
	if err := session.Save(); err != nil {
		log.Error("Failed to save session login value: %v", err)
	}
	c.Redirect(http.StatusFound, "/")
}

func (w *Web) getLogin(c *gin.Context) {
	w.render(c, login, w.defaultM(c, login))
}

func (w *Web) getFireTracker(c *gin.Context) {
	fw, err := apiGetFireWatches(w.ctx, w.client, 10)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	m := w.defaultM(c, firetracker)
	m["firewatces"] = fw
	w.render(c, firetracker, m)
}

func (w *Web) getExample(c *gin.Context) {
	idStr := c.Param("example_id")
	if idStr == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ex, err := apiGetExample(w.ctx, w.client, int(id))
	if err != nil {
		log.Errorf("Failed to fetch examples: %v", err)
	}
	m := w.defaultM(c, example)
	m["example"] = ex
	w.render(c, example, m)
}

func (w *Web) getExamples(c *gin.Context) {
	exs, err := apiGetExamples(w.ctx, w.client)
	if err != nil {
		log.Errorf("Failed to fetch examples: %v", err)
	}
	m := w.defaultM(c, examples)
	m["examples"] = exs
	w.render(c, examples, m)
}

func (w *Web) getPartners(c *gin.Context) {
	type Partner struct {
		Name string
		Desc string
		URL  string
		Img  string
	}
	var partnerBlocks = []Partner{
		{
			"Drone America",
			`Designs and builds small and large drones capable of 
						VTOL launch, 6hrs range and BVLOS up to 350lbs. These will 
						carry VeriMap geomatics mappers to provide Real-Time mapping.`,
			"droneamerica.com",
			"drone-america-partner-1024x683.png"},
		{
			"Inertial-SPI (a division of Special Projects Inc.)",
			`Team of experts in engineering, software development, visionary application 
					architecture, GIS programming`,
			"www.inertial.ca",
			"special-projects-partner-1024x683.png"},
		{
			"MARC Inc.",
			`North America's largest provider of specialized contract aircraft 
						and flight crews for airborne GIS and survey`,
			"www.marcflightservices.com",
			"marcinc-project-1024x683.png",
		}}
	m := w.defaultM(c, partners)
	m["partners"] = partnerBlocks
	w.render(c, partners, m)
}

func (w *Web) getHome(c *gin.Context) {
	type ReadMore struct {
		Title string
		Icon  string
		Desc  string
		URL   string
		Key   string
	}
	var readMoreBlocks = []ReadMore{
		{"Emergency Response",
			"fi-first-aid",
			`Speed and accuracy are crucial for
        mitigating natural disasters, saving lives, and minimizing impact to property. Verimap
        has years of experience delivering critical data fast.`,
			"",
			"emergency",
		},
		{
			"Environmental",
			"fi-trees",
			`Wide area data collection is key for modelling
        and tracking changes in Earth’s climate and ecosystems. Verimap has experience mapping large areas fast,
        and in multiple spectrums.`,
			"/services/environmental",
			"environmental",
		}, {
			"Infrastructure",
			"fi-safety-cone",
			`Monitoring and management of large
        terrestrial assets can be costly with time consuming field visits to difficult-to-access and remote
        locations. Verimap provides real asset information from the air allowing decisions before deploying
        field teams, reducing the number of field visits and creating a lasting record of assets
        over time.`,
			"/services/infrastructure",
			"infrastructure",
		},
	}
	showcases, err := apiGetShowcases(w.ctx, w.client)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	m := w.defaultM(c, home)
	m["showcases"] = showcases
	m["read_more"] = readMoreBlocks
	w.render(c, home, m)
}

func (w *Web) getProfile(c *gin.Context) {
	m := w.defaultM(c, profile)
	w.render(c, profile, m)
}

func (w *Web) updateUserFromForm(c *gin.Context, p *store.Person, emptyPasswordOk bool) bool {
	p.FirstName = c.PostForm("first_name")
	if p.FirstName == "" {
		abortFlash(c, "All fields are required", w.route(adminPeople))
		return false
	}
	p.LastName = c.PostForm("last_name")
	if p.LastName == "" {
		abortFlash(c, "All fields are required", w.route(adminPeople))
		return false
	}
	p.Email = c.PostForm("email")
	if p.Email == "" {
		abortFlash(c, "All fields are required", w.route(adminPeople))
		return false
	}
	id, err := strconv.ParseInt(c.PostForm("agency_id"), 10, 64)
	if err != nil {
		abortFlashErr(c, "Invalid Agency ID", w.route(adminPeople), err)
		return false
	}
	p.AgencyID = int(id)
	if p.AgencyID == 0 {
		abortFlash(c, "All fields are required", w.route(adminPeople))
		return false
	}
	pw := c.PostForm("password")
	pwV := c.PostForm("password_verify")
	if pw != "" && pwV != "" {
		hash, err := HashPassword(pw)
		if err != nil {
			abortFlashErr(c, "Passwords do not match", w.route(adminPeople), err)
			return false
		}
		p.PasswordHash = hash
	} else if pw != "" || pwV != "" {
		abortFlash(c, "Passwords must match", w.route(adminPeople))
		return false
	} else if pw == "" && pwV == "" && !emptyPasswordOk {
		abortFlash(c, "Passwords cannot be empty", w.route(adminPeople))
		return false
	}
	if err := store.SavePerson(w.ctx, w.db, p); err != nil {
		abortFlashErr(c, "Failed to save user", w.route(adminPeopleEdit), err)
		return false
	}
	return true
}

func (w *Web) postProfile(c *gin.Context) {
	m := w.defaultM(c, profile)
	var p store.Person
	if err := store.LoadPersonByID(w.ctx, w.db, m["person"].(store.Person).PersonID, &p); err != nil {
		abortFlashErr(c, "Error loading person data", w.route(profile), err)
		return
	}
	if w.updateUserFromForm(c, &p, true) {
		successFlash(c, "Updated profile successfully", w.route(profile))
	}
}
