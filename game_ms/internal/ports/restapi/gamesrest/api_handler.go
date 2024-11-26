package gamesrest

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"dataxo-backend-game-ms/internal/ports/restapi"
	"github.com/google/uuid"
	"log/slog"
)

type GameUsecase interface {
	CreateGame(ctx context.Context, player domain.Player, mode string, params domain.ModeParams) (*domain.Game, error)
	GetGame(ctx context.Context, gameID uuid.UUID) (*domain.Game, error)
	MakeMove(ctx context.Context, gameID uuid.UUID, move domain.Move, side int) error
}

type Handler struct {
	log       *slog.Logger
	gameUC    GameUsecase
	responder restapi.Responder
}

func New(log *slog.Logger, gameUC GameUsecase, responder restapi.Responder) *Handler {
	return &Handler{log: log, gameUC: gameUC, responder: responder}
}
