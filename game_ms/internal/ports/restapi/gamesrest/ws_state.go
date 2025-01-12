package gamesrest

import (
	"dataxo-backend-game-ms/internal/domain"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	"log/slog"
)

type WsGameConfig struct {
	// 0 is no limit
	PlayerFiguresLimit int `json:"player_figures_limit"`
	WinLineLength      int `json:"win_line_length"`
	BoardWidth         int `json:"board_width"`
	BoardHeight        int `json:"board_height"`
}

type WsGameStateResp struct {
	Type          string `json:"type"`
	ResponseForID string `json:"response_for_id"`

	GameID      uuid.UUID      `json:"game_id"`
	Mode        string         `json:"mode"`
	Config      WsGameConfig   `json:"config"`
	State       domain.State   `json:"state"`
	Moves       []domain.Move  `json:"moves"`
	WinSequence []domain.Move  `json:"win_sequence"`
	Winner      domain.WinSide `json:"winner"`
}

var GameStateResponseType = "game_state_response"

func (h *Handler) WsState(session *melody.Session, requestID string, gameID uuid.UUID, bytes []byte) {
	ctx := session.Request.Context()

	g, err := h.gameUC.GetGame(ctx, gameID)
	if err != nil {
		h.log.Error("ws state: get game", slog.Any("error", err))
		h.WsRespondErrorWithID(session, err, requestID)
		return
	}

	cfg := g.Config

	h.wsResponder.RespondWs(session, WsGameStateResp{
		Type:          GameStateResponseType,
		ResponseForID: requestID,
		GameID:        gameID,
		Mode:          g.Mode,
		Config: WsGameConfig{
			PlayerFiguresLimit: cfg.PlayerFiguresLimit,
			WinLineLength:      cfg.WinLineLength,
			BoardWidth:         cfg.BoardWidth,
			BoardHeight:        cfg.BoardHeight,
		},
		State:       g.State,
		Moves:       g.Moves,
		WinSequence: g.WinSequence,
		Winner:      g.Winner,
	})
}
