package gamesrest

import (
	"dataxo-backend-game-ms/internal/domain"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	"log/slog"
)

type WsGameReq struct {
	MoveID int `json:"move_id"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type WsGameResp struct {
	Type          string `json:"type"`
	ResponseForID string `json:"response_for_id"`
}

type MoveEvent struct {
	Type      string      `json:"type"`
	MoveID    int         `json:"move_id"`
	X         int         `json:"x"`
	Y         int         `json:"y"`
	Side      domain.Side `json:"side"`
	TimesUsed int         `json:"times_used"`
}

func (e *MoveEvent) FromDomain(event domain.MoveEvent) {
	t := ""
	switch event.Type {
	case domain.PlaceMove:
		t = "place_move"
	case domain.RemoveMove:
		t = "remove_move"
	case domain.HeatMove:
		t = "heat_move"
	case domain.BlockMove:
		t = "block_move"
	}

	move := event.Move

	e.Type = t
	e.MoveID = move.InGameID
	e.X = move.X
	e.Y = move.Y
	e.Side = move.Side
	e.TimesUsed = move.TimesUsed
}

func MoveEventsFromDomain(events []domain.MoveEvent) []MoveEvent {
	converted := make([]MoveEvent, len(events))

	for i := range events {
		converted[i].FromDomain(events[i])
	}

	return converted
}

type WsMoveBroadcast struct {
	Type       string      `json:"type"`
	MoveEvents []MoveEvent `json:"move_events"`
}

type WsGameFinishBroadcast struct {
	Type string `json:"type"`
}

var (
	SuccessMoveResponseType = "success_move_response"
	MoveBroadcastType       = "new_move_broadcast"
	GameFinishBroadcastType = "game_finish_broadcast"
)

func (h *Handler) WsGame(session *melody.Session, requestID string, gameID uuid.UUID, bytes []byte) {
	ctx := session.Request.Context()

	req := &WsGameReq{}
	if err := json.Unmarshal(bytes, req); err != nil {
		h.WsRespondErrorWithID(session, err, requestID)
		return
	}

	h.log.Debug("WebSocket Game Request",
		slog.String("json", string(bytes)),
		slog.Any("struct", req),
	)

	side, err := h.WsGetSide(ctx, session, gameID)
	if err != nil {
		h.log.Error("get side", slog.Any("error", err))
		h.WsRespondErrorWithID(session, err, requestID)
		return
	}

	res, err := h.gameUC.MakeMove(ctx, gameID, domain.Move{
		InGameID: req.MoveID,
		X:        req.X,
		Y:        req.Y,
		Side:     side,
	})
	if err != nil {
		h.WsRespondErrorWithID(session, err, requestID)
		return
	}

	/*
		h.wsResponder.RespondWs(session, WsGameResp{
			Type:          SuccessMoveResponseType,
			ResponseForID: requestID,
		})
	*/

	data, _ := h.wsResponder.Marshal(WsMoveBroadcast{
		Type:       MoveBroadcastType,
		MoveEvents: MoveEventsFromDomain(res.Events),
	})

	err = h.wsHandler.BroadcastFilter(data, func(other *melody.Session) bool {
		otherGameID, err := h.WsGameIDFromSession(other)
		if err != nil {
			h.log.Error("game move broadcast", slog.Any("error", err))
		}

		if gameID == otherGameID {
			return true
		}
		return false
	})
	if err != nil {
		h.log.Error("ws move broadcast", slog.Any("error", err))
	}

	if !res.GameFinished {
		return
	}

	data, _ = h.wsResponder.Marshal(WsGameFinishBroadcast{
		Type: GameFinishBroadcastType,
	})

	err = h.wsHandler.Broadcast(data)
	if err != nil {
		h.log.Error("ws game finish broadcast", slog.Any("error", err))
	}
}
