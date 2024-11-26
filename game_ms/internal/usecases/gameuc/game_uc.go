package gameuc

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"github.com/google/uuid"
)

type GameRepository interface {
	CreateGame(ctx context.Context, player domain.Player, mode string, params domain.ModeParams) (*domain.Game, error)
	GetGame(ctx context.Context, gameID uuid.UUID) (*domain.Game, error)
}

type GameUC struct {
	gameRepo GameRepository
}

func (uc *GameUC) CreateGame(ctx context.Context, pl domain.Player, mode string, params domain.ModeParams) (*domain.Game, error) {
	return uc.gameRepo.CreateGame(ctx, pl, mode, params)
}

func (uc *GameUC) GetGame(ctx context.Context, gameID uuid.UUID) (*domain.Game, error) {
	return uc.gameRepo.GetGame(ctx, gameID)
}

func (uc *GameUC) MakeMove(ctx context.Context, gameID uuid.UUID, move domain.Move, side int) error {
	return nil
}
