package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/leighmacdonald/verimapcom/web/store"
	"net/http"
	"strconv"
	"time"
)

func (w *Web) postMission(c *gin.Context) {
	var mission store.Mission
	m := w.defaultM(c, missionsCreate)
	u := m["person"].(store.Person)
	mission.AgencyID = u.AgencyID
	mission.MissionState = int(store.StateCreated)
	mission.PersonID = u.PersonID
	mission.CreatedOn = time.Now()
	mission.UpdatedOn = mission.CreatedOn
	mission.BoundingBox.LatUL = formFloatDefault(c, "lat_ul", 0)
	mission.BoundingBox.LongUL = formFloatDefault(c, "lon_ul", 0)
	mission.BoundingBox.LatLR = formFloatDefault(c, "lat_lr", 0)
	mission.BoundingBox.LongLR = formFloatDefault(c, "lon_lr", 0)
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

func (w *Web) getMission(c *gin.Context) {
	missionID, err := strconv.ParseInt(c.Param("mission_id"), 10, 64)
	if err != nil {
		abortFlash(c, "Invalid mission", referer(c))
		return
	}
	mis, err := store.GetMission(w.ctx, w.db, int(missionID))
	if err != nil {
		abortFlashErr(c, "Failed to load mission", referer(c), err)
		return
	}
	files, err := store.FileGetAllMission(w.ctx, w.db, mis.MissionID)
	if err != nil {
		abortFlashErr(c, "Failed to load mission files", referer(c), err)
		return
	}
	flights, err := store.FlightsByMissionID(w.ctx, w.db, mis.MissionID)
	if err != nil {
		abortFlashErr(c, "Failed to load flights", referer(c), err)
		return
	}
	m := w.defaultM(c, mission)
	m["mission"] = mis
	m["files"] = files
	m["flights"] = flights
	w.render(c, mission, m)
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

func (w *Web) getMissionEvents(c *gin.Context) {
	missionID, err := strconv.ParseInt(c.Param("mission_id"), 10, 64)
	if err != nil {
		abortFlash(c, "Invalid mission", referer(c))
		return
	}
	mis, err := store.GetMission(w.ctx, w.db, int(missionID))
	if err != nil {
		abortFlashErr(c, "Failed to load mission", referer(c), err)
		return
	}
	events, err := store.MissionEventGetAll(w.ctx, w.db, mis.MissionID)
	if err != nil {
		abortFlashErr(c, "Failed to load events", referer(c), err)
		return
	}
	c.JSON(http.StatusOK, events)
}
