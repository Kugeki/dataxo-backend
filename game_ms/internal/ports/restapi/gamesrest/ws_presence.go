package gamesrest

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	"log/slog"
)

var (
	ErrCantGetGameIDFromSession     = errors.New("can't get game id from session")
	ErrCantConvertGameIDValueToUUID = errors.New("can't convert game id value to uuid")
)

type WsPresenceReq struct {
	Action string `json:"action"`
}

type WsGameStartBroadcast struct {
	Type string `json:"type"`
}

var GameStartBroadcastType = "start_broadcast"

func (h *Handler) WsPresence(session *melody.Session, requestID string, gameID uuid.UUID, bytes []byte) {
	ctx := session.Request.Context()

	req := &WsPresenceReq{}
	if err := json.Unmarshal(bytes, req); err != nil {
		h.WsRespondErrorWithID(session, err, requestID)
		return
	}

	h.log.Debug("WebSocket Presence Request",
		slog.String("json", string(bytes)),
		slog.Any("struct", req),
	)

	gameID, err := h.WsGameIDFromSession(session)
	if err != nil {
		h.WsRespondErrorWithID(session, err, requestID)
		return
	}

	playerID := h.WsGetPlayerID(session)

	switch req.Action {
	case "join":
		res, err := h.gameUC.JoinGame(ctx, gameID, playerID)
		if err != nil {
			h.WsRespondErrorWithID(session, err, requestID)
			return
		}

		session.Set("side", res.Side)

		h.WsSendSide(ctx, session, requestID, gameID)

		if !res.ReadyToStart {
			return
		}

		err = h.gameUC.StartGame(ctx, gameID)
		if err != nil {
			h.WsRespondErrorWithID(session, err, requestID)
			return
		}

		data, _ := h.wsResponder.Marshal(WsGameStartBroadcast{Type: GameStartBroadcastType})

		err = h.wsHandler.Broadcast(data)
		if err != nil {
			h.log.Error("ws handler broadcast", slog.Any("error", err))
		}
	case "leave":
		// todo
	default:
		h.WsRespondErrorWithID(session, &PresenceActionError{
			Err:    ErrInvalidPresenceAction,
			Action: req.Action}, requestID)
		return
	}

	// todo
}
