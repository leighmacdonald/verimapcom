package client

import (
	"github.com/hpcloud/tail"
	"strings"
	"time"
)

type TimedLine struct {
	Text string
	Time time.Time
}

func watchFile(t *tail.Tail, outChan chan TimedLine) {
	for line := range t.Lines {
		tl := TimedLine{strings.ReplaceAll(line.Text, "\r", ""), line.Time}
		outChan <- tl
	}
}
