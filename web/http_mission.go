package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/leighmacdonald/verimapcom/web/store"
	"time"
)

func (w *Web) postMission(c *gin.Context) {
	var mission store.Mission
	m := w.defaultM(c, missionsCreate)
	u := m["person"].(store.Person)
	mission.AgencyID = u.AgencyID
	mission.MissionState = store.StateCreated
	mission.PersonID = u.PersonID
	mission.CreatedOn = time.Now()
	mission.UpdatedOn = mission.CreatedOn
	name := c.PostForm("mission_name")
	if name == "" {
		abortFlash(c, "Invalid mission name, cannot be empty", w.route(missionsCreate))
		return
	}
	mission.MissionName = name
	if err := store.SaveMission(w.ctx, w.db, &mission); err != nil {
		abortFlashErr(c, "Failed to save mission", w.route(missionsCreate), err)
		return
	}
	successFlash(c, "Create mission successfully", w.route(missions))
}

func (w *Web) getMissionsCreate(c *gin.Context) {
	m := w.defaultM(c, missionsCreate)
	agencies, err := store.GetAgencies(w.ctx, w.db)
	if err != nil {
		abortFlashErr(c, "Failed to get agencies", w.route(home), err)
		return
	}
	m["agencies"] = agencies
	w.render(c, missionsCreate, m)
}

func (w *Web) getMissions(c *gin.Context) {
	m := w.defaultM(c, missions)
	userMissions, err := store.GetMissions(w.ctx, w.db, m["person"].(store.Person).AgencyID)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		abortFlashErr(c, "Failed to get missions", w.route(home), err)
		return
	}
	m["missions"] = userMissions
	w.render(c, missions, m)
}
