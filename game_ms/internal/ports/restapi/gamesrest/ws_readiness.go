package gamesrest

import (
	"encoding/json"
	"github.com/olahol/melody"
	"log/slog"
)

type WsReadinessReq struct {
	Action string `json:"action"`
}

func (h *Handler) WsReadiness(session *melody.Session, bytes []byte) {
	// todo
	req := &WsReadinessReq{}
	if err := json.Unmarshal(bytes, req); err != nil {
		h.wsResponder.RespondErrorWs(session, err)
		return
	}

	h.log.Debug("WebSocket Readiness Request",
		slog.String("json", string(bytes)),
		slog.Any("struct", req),
	)

	switch req.Action {
	case "ready":
		// todo
	case "unready":
		// todo
	default:
		h.wsResponder.RespondErrorWs(session, &ReadinessError{
			Err: ErrInvalidReadinessAction, Action: req.Action,
		})
		return
	}

	// todo
}
