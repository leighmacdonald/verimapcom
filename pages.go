package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func simple(page pageName) gin.HandlerFunc {
	return func(c *gin.Context) {
		render(c, page, defaultM(c, page))
	}
}

func getLogout(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("person_id")
	if v == nil {
		c.Redirect(http.StatusFound, pages[login].Path)
		return
	}
	_, ok := v.(int)
	if !ok {
		c.Redirect(http.StatusFound, pages[login].Path)
		return
	}
	session.Clear()
	if err := session.Save(); err != nil {
		log.Errorf("failed to clear user session on logout: %v", err)
	}
	c.Redirect(http.StatusFound, "/")
}

func postLogin(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("person_id")
	if v != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	email := c.PostForm("email")
	password := c.PostForm("password")
	var p Person
	if err := loadPersonByEmail(email, &p); err != nil {
		session.AddFlash("Invalid username or password")
		c.Redirect(http.StatusFound, pages[login].Path)
		return
	}
	if !CheckPasswordHash(password, p.PasswordHash) {
		session.AddFlash("Invalid username or password")
		c.Redirect(http.StatusFound, pages[login].Path)
		return
	}
	session.Set("person_id", p.PersonID)
	session.AddFlash("Logged in successfully")
	if err := session.Save(); err != nil {
		log.Error("Failed to save session login value: %v", err)
	}
	c.Redirect(http.StatusFound, "/")
}

func getLogin(c *gin.Context) {
	render(c, login, defaultM(c, login))
}

func getFireTracker(c *gin.Context) {
	fw, err := apiGetFireWatches(10)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	m := defaultM(c, firetracker)
	m["firewatces"] = fw
	render(c, firetracker, m)
}

func getExample(c *gin.Context) {
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
	ex, err := apiGetExample(int(id))
	if err != nil {
		log.Errorf("Failed to fetch examples: %v", err)
	}
	m := defaultM(c, example)
	m["example"] = ex
	render(c, example, m)
}

func getExamples(c *gin.Context) {
	exs, err := apiGetExamples()
	if err != nil {
		log.Errorf("Failed to fetch examples: %v", err)
	}
	m := defaultM(c, examples)
	m["examples"] = exs
	render(c, examples, m)
}

func getPartners(c *gin.Context) {
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
	m := defaultM(c, partners)
	m["partners"] = partnerBlocks
	render(c, partners, m)
}

func getHome(c *gin.Context) {
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
	showcases, err := apiGetShowcases()
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	m := defaultM(c, home)
	m["showcases"] = showcases
	m["read_more"] = readMoreBlocks
	render(c, home, m)
}

func getAdminAgencies(c *gin.Context) {
	agencies, err := getAgencies()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	m := defaultM(c, adminAgencies)
	m["agencies"] = agencies
	render(c, adminAgencies, m)
}

func getAdminMissions(c *gin.Context) {
	person, found := c.Get("person")
	if !found {
		c.Redirect(http.StatusUnauthorized, "/login")
		return
	}
	missions, err := dbGetMissions(person.(Person).AgencyID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	m := defaultM(c, adminMissions)
	m["missions"] = missions
	render(c, adminMissions, m)
}

func getAdminPeople(c *gin.Context) {
	people, err := getPeople()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	m := defaultM(c, adminPeople)
	m["people"] = people
	render(c, adminPeople, m)
}
