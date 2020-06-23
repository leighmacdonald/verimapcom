package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/leighmacdonald/verimapcom/web/store"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (

	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsClient struct {
	conn        *websocket.Conn
	missionID   int
	payloadChan chan payloadSend
}

func (ws *wsClient) send(p payloadSend) {
	ws.payloadChan <- p
}

type wsEvent int

const (
	evtConnect    wsEvent = 1
	evtPing       wsEvent = 2
	evtPong       wsEvent = 3
	evtMessage    wsEvent = 10
	evtSetMission wsEvent = 20
	evtError      wsEvent = 10000
)

func removeWSClientIndex(s []*wsClient, index int) []*wsClient {
	return append(s[:index], s[index+1:]...)
}

func (w *Web) serveWs(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("Failed to setup websocket connection: %v", err)
		return
	}
	closed := false
	defer func() {
		if !closed {
			if err := ws.Close(); err != nil {
				log.Errorf("Error closing client WS: %v", err)
			}
		}
	}()
	wsClientConn := &wsClient{
		conn:        ws,
		payloadChan: make(chan payloadSend),
	}
	w.wsClientMu.Lock()
	w.wsClient[ws] = wsClientConn
	w.wsClientMu.Unlock()
	defer func() {
		w.wsClientMu.Lock()
		defer w.wsClientMu.Unlock()
		delete(w.wsClient, ws)
		for i, v := range w.wsMissionConns {
			for j, c := range v {
				if c.conn == ws {
					w.wsMissionConns[i] = removeWSClientIndex(w.wsMissionConns[i], j)
					log.Debugf("Removed ws client from mission room %d", i)
					return
				}
			}
		}
	}()
	go w.wsWriter(wsClientConn)
	w.wsReader(wsClientConn, &closed)

}

type wsEventHandler func(*wsClient, payloadRecv) error

func (w *Web) wsJoinMission(ws *wsClient, missionID int) error {
	w.wsClientMu.RLock()
	clients, found := w.wsMissionConns[missionID]
	w.wsClientMu.RUnlock()
	if !found {
		w.wsClientMu.Lock()
		w.wsMissionConns[missionID] = []*wsClient{ws}
		ws.missionID = missionID
		w.wsClientMu.Unlock()
		log.Infof("Client joined mission: %d", missionID)
		return nil
	}
	for _, c := range clients {
		if c == ws {
			log.Warnf("Client already in mission room")
			return nil
		}
	}
	w.wsClientMu.Lock()
	w.wsMissionConns[missionID] = append(w.wsMissionConns[missionID], ws)
	w.wsClientMu.Unlock()
	log.Infof("Client joined mission: %d", missionID)
	return nil
}

func (w *Web) wsOnMessage(ws *wsClient, p payloadRecv) error {
	log.Debug(p.Payload)
	e := store.MissionEvent{}
	if err := store.MissionEventAdd(w.ctx, w.db, &e); err != nil {
		return err
	}
	return nil
}

func (w *Web) wsOnSetMission(ws *wsClient, p payloadRecv) error {
	log.Debug(p.Payload)
	missionID, ok := p.Payload["mission_id"].(float64)
	if !ok {
		return errors.Errorf("Invalid missionID")
	}
	return w.wsJoinMission(ws, int(missionID))
}

func (w *Web) wsOnPing(ws *wsClient, p payloadRecv) error {
	data, ok := p.Payload["data"].(string)
	if !ok {
		return errors.Errorf("Invalid data")
	}
	ws.send(payloadSend{
		Event: evtPong,
		Payload: map[string]string{
			"data": data,
		},
	})
	return nil
}

func (w *Web) handleWSMessage(ws *wsClient, m payloadRecv) error {
	fn, found := w.wsHandler[m.Event]
	if !found {
		return errors.New("invalid event")
	}
	return fn(ws, m)
}

func (w *Web) wsReader(ws *wsClient, closed *bool) {
	defer func() {
		if err := ws.conn.Close(); err != nil {
			log.Errorf("failed to cleanly shutdown ws reader: %v", err)
		}
	}()
	ws.conn.SetReadLimit(512)
	if err := ws.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Errorf("failed to set read deadline: %v", err)
	}
	ws.conn.SetPongHandler(func(string) error {
		if err := ws.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Errorf("failed to set read deadline: %v", err)
		}
		return nil
	})

	req := 0
	for {
		var payload payloadRecv
		if err := ws.conn.ReadJSON(&payload); err != nil {
			*closed = true
			if e, ok := err.(*websocket.CloseError); ok {
				if e.Code <= 1001 {
					return
				}
				log.Errorf("Websocket error: %v", e.Error())
				return
			}
			log.Errorf("Failed to read user payloadRecv: %v", err)
			return
		}
		if err := w.handleWSMessage(ws, payload); err != nil {
			log.Errorf("Failed to handle event payloadRecv: %v", err)
		}
		req++
	}
}

func (w *Web) wsWriter(ws *wsClient) {
	pingTicker := time.NewTicker(time.Second * 60)
	defer func() {
		pingTicker.Stop()
		if err := ws.conn.Close(); err != nil {
			log.Errorf("Error closing ws connection: %v", err)
		}
	}()
	for {
		select {
		case p := <-ws.payloadChan:
			if err := ws.conn.WriteJSON(&p); err != nil {
				log.Errorf("failed to write ws json payload: %v", err)
				continue
			}
		case <-pingTicker.C:
			if err := ws.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Errorf("Failed to set write deadline: %v", err)
			}
			if err := ws.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
