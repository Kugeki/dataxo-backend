package restapi

import (
	"github.com/olahol/melody"
	"net/http"
)

type Responder interface {
	Respond(w http.ResponseWriter, code int, data interface{})
	RespondError(w http.ResponseWriter, code int, err error)
}

type WsResponder interface {
	Marshal(data interface{}) ([]byte, error)
	RespondWs(session *melody.Session, data interface{})
	RespondErrorWs(session *melody.Session, err error)
}
