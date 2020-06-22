package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/leighmacdonald/verimapcom/web/store"
	"net/http"
	"strconv"
	"time"
)

func (w *Web) getAdminAgencies(c *gin.Context) {
	agencies, err := store.GetAgencies(w.ctx, w.db)
	if err != nil {
		abortFlashErr(c, "Error loading agencies", w.route(adminAgencies), err)
		return
	}
	m := w.defaultM(c, adminAgencies)
	m["agencies"] = agencies
	w.render(c, adminAgencies, m)
}

func (w *Web) postAdminAgenciesCreate(c *gin.Context) {
	name := c.PostForm("agency_name")
	slotsStr := c.PostForm("agency_slots")
	slots, err := strconv.ParseInt(slotsStr, 10, 64)
	if err != nil {
		abortFlash(c, fmt.Sprintf("Invalid slot value: %s", slotsStr), w.route(adminAgencies))
		return
	}
	if slots <= 0 {
		abortFlash(c, fmt.Sprintf("Invalid slot value, must be greater than 0 %s", slotsStr), w.route(adminAgencies))
		return
	}
	if name == "" {
		abortFlash(c, "Name cannot be empty", w.route(adminAgencies))
		return
	}
	a := store.Agency{
		AgencyName: name,
		Slots:      int(slots),
		CreatedOn:  time.Now(),
	}
	if err := store.SaveAgency(w.ctx, w.db, &a); err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			abortFlash(c, "Name already taken", w.route(adminAgencies))
			return
		}
		abortFlashErr(c, "Failed to save new agency", w.route(adminAgencies), err)
		return
	}
	successFlash(c, fmt.Sprintf("Created agency successfully: %s", name), w.route(adminAgencies))
}

func (w *Web) getAdminPeople(c *gin.Context) {
	people, err := store.GetPeople(w.ctx, w.db)
	if err != nil {
		abortFlashErr(c, "Failed to load people", w.route(home), err)
		return
	}
	agencies, err := store.GetAgencies(w.ctx, w.db)
	if err != nil {
		abortFlashErr(c, "Failed to load agencies", w.route(home), err)
		return
	}
	m := w.defaultM(c, adminPeople)
	m["people"] = people
	m["agencies"] = agencies
	w.render(c, adminPeople, m)
}

func (w *Web) getAdminPeopleEdit(c *gin.Context) {
	var person store.Person
	p, err := strconv.ParseInt(c.Param("person_id"), 10, 64)
	if err != nil {
		abortFlashErr(c, "Could not parse person_id", w.route(adminPeople), err)
		return
	}
	if err := store.LoadPersonByID(w.ctx, w.db, int(p), &person); err != nil {
		abortFlashErr(c, "Failed to load user data", w.route(adminPeople), err)
		return
	}

	agencies, err := store.GetAgencies(w.ctx, w.db)
	if err != nil {
		abortFlashErr(c, "Failed to load agencies", w.route(adminPeople), err)
		return
	}
	m := w.defaultM(c, adminPeopleEdit)
	m["agencies"] = agencies
	m["subject"] = person
	w.render(c, adminPeopleEdit, m)
}

func (w *Web) postAdminPeopleCreate(c *gin.Context) {
	if w.updateUserFromForm(c, &store.Person{}, true) {
		successFlash(c, "Created user successfully", w.route(adminPeople))
	}
}

func (w *Web) postAdminPeopleDelete(c *gin.Context) {
	pid, err := strconv.ParseInt(c.Param("person_id"), 10, 64)
	if err != nil {
		abortFlash(c, "Failed to find user", w.route(adminPeople))
		return
	}
	var p store.Person
	if err := store.LoadPersonByID(w.ctx, w.db, int(pid), &p); err != nil {
		abortFlash(c, "Failed to find user", w.route(adminPeople))
		return
	}
	if err := store.PersonDelete(w.ctx, w.db, p.PersonID); err != nil {
		abortFlashErr(c, "Failed to delete user", w.route(adminPeople), err)
		return
	}
	successFlash(c, "Deleted user successfully", w.route(adminPeople))
}

func (w *Web) postAdminPeopleEdit(c *gin.Context) {
	pid, err := strconv.ParseInt(c.Param("person_id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var p store.Person
	if err := store.LoadPersonByID(w.ctx, w.db, int(pid), &p); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if w.updateUserFromForm(c, &p, true) {
		successFlash(c, "Updated user successfully", w.route(adminPeople))
	}
}
