package web

import "github.com/leighmacdonald/verimapcom/web/store"

type payloadRecv struct {
	Event   store.Evt              `json:"event"`
	Payload map[string]interface{} `json:"payload"`
}

type payloadSend struct {
	Event   store.Evt   `json:"event"`
	Payload interface{} `json:"payload"`
}

type payloadMessage struct {
	MissionID  int    `json:"mission_id"`
	PersonName string `json:"person_name"`
	PersonID   int    `json:"person_id"`
	Message    string `json:"message"`
}

type payloadError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
