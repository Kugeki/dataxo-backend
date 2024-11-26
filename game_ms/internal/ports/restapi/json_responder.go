package restapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type JsonResponder struct {
	log *slog.Logger
}

func (r *JsonResponder) Respond(w http.ResponseWriter, code int, data interface{}) {
	log := r.log

	var jsonData []byte
	var err error

	if data == nil {
		log.Warn("respond helper: data is nil")
	}

	jsonData, err = json.Marshal(data)

	if err != nil {
		log.Error("respond helper: json marshal error", slog.Any("error", err))

		jsonData = []byte("{\"error\": \"json marshal error\"}")
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(code)
	_, err = w.Write(jsonData)
	if err != nil {
		log.Error("respond helper: response write json data error", slog.Any("error", err))
	}
}

func (r *JsonResponder) RespondError(w http.ResponseWriter, code int, err error) {
	r.Respond(w, code, &HTTPError{Error: err.Error()})
}
