package gamesrest

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"github.com/google/uuid"
)

type GameUsecase interface {
	CreateGame(ctx context.Context, playerID domain.PlayerID, mode string, params domain.ModeParams) (*domain.Game, error)
	GetGame(ctx context.Context, gameID uuid.UUID) (*domain.Game, error)
	JoinGame(ctx context.Context, gameID uuid.UUID, playerID domain.PlayerID) (domain.JoinGameResult, error)
	StartGame(ctx context.Context, gameID uuid.UUID) error
	MakeMove(ctx context.Context, gameID uuid.UUID, move domain.Move) (domain.MakeMoveResult, error)
	GetSide(ctx context.Context, gameID uuid.UUID, playerID domain.PlayerID) (domain.Side, error)
}
