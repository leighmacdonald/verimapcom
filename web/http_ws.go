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
	evtMissionOpen wsEvent = iota + 1
	evtMissionEvents
	evtMissionSendMessage
	evtMissionRecvMessage
	evtMissionNewFlight
	evtMissionPosition
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

type wsPayloadMissionOpen struct {
	MissionID int32
}

type wsPayloadSendMessage struct {
	Message string
}

func (ep *wsEventPayload) Decode(i interface{}) error {
	if err := json.NewDecoder(bytes.NewBuffer(ep.Payload)).Decode(i); err != nil {
		return errors.Wrapf(err, "failed to decode to interface")
	}
	return nil
}

func (s *wsSession) reader(w *Web) {
	for {
		var err error
		var payload wsEventPayload
		err = s.conn.ReadJSON(&payload)
		if err != nil {
			switch e := err.(type) {
			case *websocket.CloseError:
				//log.Errorf("Client Disconnected: %e", e)
			default:
				log.Errorf("Failed to read JSON payload: %v", e)
			}
			break
		}
		switch payload.Event {
		case evtMissionOpen:
			var p wsPayloadMissionOpen
			if err := json.Unmarshal(payload.Payload, &p); err != nil {

			}
			err = w.wsHandleMissionOpen(p)
		case evtMissionEvents:
			err = w.wsHandleMissionEvents()
		case evtMissionNewFlight:
			err = w.wsHandleMissionNewFlight()
		case evtMissionRecvMessage:
			err = w.wsHandleMissionRecvMessage()
		case evtMissionSendMessage:
			var p wsPayloadSendMessage
			if err2 := json.Unmarshal(payload.Payload, &p); err2 != nil {
				err = err2
				break
			}
			err = w.wsHandleMissionSendMessage(p)
		case evtMissionPosition:
			err = w.wsHandlePosition()
		}
		if err != nil {
			log.Errorf("failed to handle ws even: %v", err)
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
	log.Infof("Client ws handler conn opened")
	c, cancel := context.WithCancel(w.ctx)
	defer cancel()
	session := newWSSession(c, conn)
	go session.writer()
	session.reader(w)
	log.Infof("Client ws handler conn closed")
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

func (w *Web) wsHandleMissionOpen(p wsPayloadMissionOpen) error {
	panic("unimplemented")
}

func (w *Web) wsHandleMissionEvents() error {
	panic("unimplemented")
}

func (w *Web) wsHandleMissionSendMessage(p wsPayloadSendMessage) error {
	log.Println(p)
	return nil
}

func (w *Web) wsHandleMissionRecvMessage() error {
	panic("unimplemented")
}

func (w *Web) wsHandleMissionNewFlight() error {
	panic("unimplemented")
}

func (w *Web) wsHandlePosition() error {
	panic("unimplemented")
}
