package web

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	log "github.com/sirupsen/logrus"
	"html/template"
	"strconv"
	"strings"
)

type M map[string]interface{}

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
	Method  string
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
