package restapi

import "net/http"

type Responder interface {
	Respond(w http.ResponseWriter, code int, data interface{})
	RespondError(w http.ResponseWriter, code int, err error)
}
