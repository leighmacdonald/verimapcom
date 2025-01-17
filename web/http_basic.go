package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leighmacdonald/verimapcom/store"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func referer(c *gin.Context) string {
	r := c.Request.Referer()
	return r
}

func (w *Web) simple(page pageName) gin.HandlerFunc {
	return func(c *gin.Context) {
		w.render(c, page, w.defaultM(c, page))
	}
}

func (w *Web) sendConnectMessage(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	subject := c.PostForm("subject")
	body := c.PostForm("body")
	if name == "" {
		abortFlash(c, "Name cannot be empty", referer(c))
		return
	}
	if email == "" {
		abortFlash(c, "Email cannot be empty", referer(c))
		return
	}
	if subject == "" {
		abortFlash(c, "Subject cannot be empty", referer(c))
		return
	}
	if body == "" {
		abortFlash(c, "Message body cannot be empty", referer(c))
		return
	}
	m := store.Message{
		Author:    email,
		AuthorID:  0,
		Subject:   fmt.Sprintf("[%s] %s", name, subject),
		Body:      body,
		CreatedOn: time.Time{},
	}
	if err := store.MessageAdd(w.ctx, w.db, &m); err != nil {
		abortFlashErr(c, "Failed to send messag, please try again later", referer(c), err)
		return
	}
	successFlash(c, "Successfully sent message, thanks for your inquiry", referer(c))
}

func (w *Web) getBackground(c *gin.Context) {
	type teamMember struct {
		Name        string
		Credentials string
		Role        string
		Image       string
		Email       string
	}
	members := []teamMember{
		{
			Name:        "David Stonehouse",
			Credentials: "Founder/CEO",
			Role:        "Certified Thermographer",
			Image:       "https://get.foundation/sites/docs/assets/img/generic/rectangle-1.jpg",
			Email:       "dstonehouse@verimap.com",
		},
		{
			Name:        "Rod Coppock",
			Credentials: "CPA, CA",
			Role:        "Chief Financial Officer",
			Image:       "https://get.foundation/sites/docs/assets/img/generic/rectangle-1.jpg",
			Email:       "",
		},
		{
			Name:        "Erkin Atakhanov",
			Credentials: "CPA, CA",
			Role:        "Corporate Finance & Investor Relations",
			Image:       "https://get.foundation/sites/docs/assets/img/generic/rectangle-1.jpg",
			Email:       "",
		}, {
			Name:        "Dale Good",
			Credentials: "P. Eng",
			Role:        "Technical Advisor",
			Image:       "https://get.foundation/sites/docs/assets/img/generic/rectangle-1.jpg",
			Email:       "",
		},
	}
	m := w.defaultM(c, background)
	m["members"] = members
	w.render(c, background, m)
}

func (w *Web) getFireTracker(c *gin.Context) {
	fw, err := apiGetFireWatches(w.ctx, w.client, 10)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	m := w.defaultM(c, firetracker)
	m["firewatches"] = fw
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
