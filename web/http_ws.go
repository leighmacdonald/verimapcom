package web

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type wsEvent int

const (
	openProject wsEvent = 1
)

type wsSession struct {
	MissionID int32
	sendChan  chan []byte
	done      chan interface{}
	conn      *websocket.Conn
	ctx       context.Context
}

func newWSSession(ctx context.Context, conn *websocket.Conn) *wsSession {
	return &wsSession{
		MissionID: 0,
		sendChan:  make(chan []byte),
		conn:      conn,
		ctx:       ctx,
	}
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsEventPayload struct {
	Event   wsEvent         `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

func (ep *wsEventPayload) Decode(i interface{}) error {
	if err := json.NewDecoder(bytes.NewBuffer(ep.Payload)).Decode(i); err != nil {
		return errors.Wrapf(err, "failed to decode to interface")
	}
	return nil
}

func (s *wsSession) reader() {
	for {
		var payload wsEventPayload
		err := s.conn.ReadJSON(&payload)
		if err != nil {
			break
		}

	}
}

func (s *wsSession) writer() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case p := <-s.sendChan:
			if err := s.conn.WriteMessage(websocket.TextMessage, p); err != nil {
				log.Errorf("Failed to send & encode full payload: %v", err)
				return
			}
		}
	}
}

func wsHandler(w *Web, ctx *gin.Context) {
	conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Errorf("Failed to set websocket upgrade: %+v", err)
		return
	}
	c, cancel := context.WithCancel(w.ctx)
	defer cancel()
	session := newWSSession(c, conn)
	go session.writer()
	session.reader()
}

func wsBroadcast(sessions []*wsSession, event wsEvent, payload interface{}) error {
	a, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrapf(err, "Failed to preencode payload")
	}
	b, err := json.Marshal(wsEventPayload{Event: event, Payload: a})
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal json broadcast payload")
	}
	for _, s := range sessions {
		if err := s.conn.WriteMessage(websocket.TextMessage, b); err != nil {
			log.Errorf("Failed to send payload to %s: %v", s.conn.RemoteAddr(), err)
		}
	}
	return nil
}

func wsSend(conn *websocket.Conn, event wsEvent, payload interface{}) error {
	a, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrapf(err, "Failed to preencode payload")
	}
	if err := conn.WriteJSON(wsEventPayload{
		Event:   event,
		Payload: a,
	}); err != nil {
		return errors.Wrapf(err, "Failed to send & encode full payload")
	}
	return nil
}
