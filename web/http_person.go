package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/leighmacdonald/verimapcom/web/store"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

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
