package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leighmacdonald/verimapcom/web/store"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	// 10g
	FileMaxBytes = 1024 * 1024 * 1024 * 1024 * 10
)

func (w *Web) postUpload(c *gin.Context) {
	missionId, err := strconv.ParseInt(c.PostForm("mission_id"), 10, 64)
	if err != nil {
		abortFlash(c, "Invalid mission_id", referer(c))
		return
	}
	fFile, err := c.FormFile("file")
	if err != nil {
		return
	}
	if fFile.Size > FileMaxBytes {
		abortFlash(c, "Filesize too large, max 10gb", referer(c))
		return
	}
	fp, err := fFile.Open()
	if err != nil {
		abortFlash(c, "Failed to process file", referer(c))
		return
	}
	defer func() {
		if err := fp.Close(); err != nil {
			log.Errorf("Failed to close uploaded file pointer: %v", err)
		}
	}()
	p, err := w.currentPerson(c)
	if err != nil {
		abortFlashErr(c, "Failed to authenticate upload", referer(c), err)
		return
	}

	f, err := store.NewFile(fp, fFile.Filename, p.PersonID)
	if err != nil {
		abortFlashErr(c, "Failed to create file", referer(c), err)
		return
	}
	if err := store.FileSave(w.ctx, w.db, w.uploadPath, f); err != nil {
		abortFlashErr(c, "Failed to save file", referer(c), err)
		return
	}
	if missionId > 0 {
		if err := store.MissionAttachFile(w.ctx, w.db, int(missionId), f.FileID); err != nil {
			abortFlash(c, "Failed to attach file to mission", referer(c))
			return
		}
	}
	successFlash(c, "Uploaded file successfully", w.route(uploads))
}

func (w *Web) getFile(c *gin.Context) {
	fid, err := strconv.ParseInt(c.Param("file_id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	file, err := store.FileGet(w.ctx, w.db, int(fid))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	p, err := w.currentPerson(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if !store.FileHaveAccess(w.ctx, w.db, p.AgencyID, file.FileID, p.PersonID) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if err := store.FileRead(w.uploadPath, file); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Description", "Verimap File Download")
	c.Header("Content-Type", file.FileType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
	c.Data(http.StatusOK, file.FileType, file.Data)

	if err := store.FileRegisterDownload(w.ctx, w.db, p.PersonID, file.FileID); err != nil {
		log.Errorf("Failed to register download: %v", err)
	}
}

func (w *Web) getUpload(c *gin.Context) {
	w.render(c, upload, w.defaultM(c, upload))
}

func (w *Web) getUploads(c *gin.Context) {
	m := w.defaultM(c, uploads)
	p := m["person"].(store.Person)
	files, err := store.FileUploadsGetPaged(w.ctx, w.db, p, 100, 0)
	if err != nil {
		abortFlashErr(c, "Failed to get downloads", w.route(home), err)
		return
	}
	m["files"] = files
	w.render(c, uploads, m)
}

func (w *Web) getDownloads(c *gin.Context) {
	m := w.defaultM(c, downloads)
	p := m["person"].(store.Person)
	files, err := store.FileGetPaged(w.ctx, w.db, p, 100, 0)
	if err != nil {
		abortFlashErr(c, "Failed to get downloads", w.route(home), err)
		return
	}
	m["files"] = files
	w.render(c, downloads, m)
}
