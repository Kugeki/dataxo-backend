package restapi

import (
	"dataxo-backend-game-ms/pkg/slogdiscard"
	"encoding/json"
	"github.com/olahol/melody"
	"log/slog"
	"net/http"
)

type JsonResponder struct {
	log *slog.Logger
}

func NewJsonResponder(log *slog.Logger) *JsonResponder {
	log = slogdiscard.LoggerIfNil(log)

	return &JsonResponder{log: log}
}

func (r *JsonResponder) Marshal(data interface{}) ([]byte, error) {
	log := r.log

	var jsonBytes []byte
	var err error

	if data == nil {
		log.Warn("respond helper: data is nil")
	}

	jsonBytes, err = json.Marshal(data)

	if err != nil {
		log.Error("respond helper: json marshal error", slog.Any("error", err))

		jsonBytes = []byte("{\"error\": \"json marshal error\"}")
		return jsonBytes, err
	}

	return jsonBytes, nil
}

func (r *JsonResponder) Respond(w http.ResponseWriter, code int, data interface{}) {
	log := r.log

	jsonBytes, err := r.Marshal(data)
	if err != nil {
		log.Error("respond helper: json marshal error", slog.Any("error", err))
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(code)
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Error("respond helper: response write json data error", slog.Any("error", err))
	}
}

func (r *JsonResponder) RespondError(w http.ResponseWriter, code int, err error) {
	r.Respond(w, code, &HTTPError{Error: err.Error()})
}

func (r *JsonResponder) RespondWs(session *melody.Session, data interface{}) {
	jsonBytes, _ := r.Marshal(data)
	r.RespondWsBytes(session, jsonBytes)
}

func (r *JsonResponder) RespondWsBytes(session *melody.Session, b []byte) {
	err := session.Write(b)
	if err != nil {
		r.log.Error("ws respond helper: response write json data error", slog.Any("error", err))
	}
}

func (r *JsonResponder) RespondErrorWs(session *melody.Session, err error) {
	r.RespondWs(session, &HTTPError{Error: err.Error()})
}
