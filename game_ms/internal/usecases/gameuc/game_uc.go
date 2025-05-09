package gameuc

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"dataxo-backend-game-ms/pkg/slogdiscard"
	"github.com/google/uuid"
	"log/slog"
	"math/rand"
)

type GameRepository interface {
	CreateGame(ctx context.Context, plID domain.PlayerID, side domain.Side, mode string, cfg domain.DisappearingModeConfig) (*domain.Game, error)
	UpdateGameState(ctx context.Context, gameID uuid.UUID, state domain.State) error
	UpdateGame(ctx context.Context, g *domain.Game) error
	GetGame(ctx context.Context, gameID uuid.UUID) (*domain.Game, error)
	GetPlayers(ctx context.Context, gameID uuid.UUID) (x *domain.Player, o *domain.Player, err error)
	AddGamePlayer(ctx context.Context, gameID uuid.UUID, playerID domain.PlayerID, side domain.Side) error
}

type GameMode interface {
	IterateGame(ctx context.Context, g *domain.Game, move domain.Move) error
	GetConfig() domain.DisappearingModeConfig
}

type GameUC struct {
	gameRepo GameRepository
	gameMode GameMode
	log      *slog.Logger
}

func New(gameRepo GameRepository, gameMode GameMode, log *slog.Logger) *GameUC {
	return &GameUC{gameRepo: gameRepo, gameMode: gameMode, log: slogdiscard.LoggerIfNil(log)}
}

func (uc *GameUC) CreateGame(ctx context.Context, plID domain.PlayerID, mode string, params domain.ModeParams) (*domain.Game, error) {
	var side domain.Side
	switch params.MySide {
	case domain.XSideRequest:
		side = domain.XSide
	case domain.OSideRequest:
		side = domain.OSide
	case domain.RandomSideRequest:
		side = domain.Side(rand.Int()%2 + 1)
	default:
		return nil, domain.ErrInvalidSide
	}

	return uc.gameRepo.CreateGame(ctx, plID, side, mode, uc.gameMode.GetConfig())
}

func (uc *GameUC) GetGame(ctx context.Context, gameID uuid.UUID) (*domain.Game, error) {
	return uc.gameRepo.GetGame(ctx, gameID)
}

func (uc *GameUC) JoinGame(ctx context.Context, gameID uuid.UUID, playerID domain.PlayerID) (domain.JoinGameResult, error) {
	xPlayer, oPlayer, err := uc.gameRepo.GetPlayers(ctx, gameID)
	if err != nil {
		return domain.JoinGameResult{}, err
	}

	if (xPlayer != nil && xPlayer.ID == playerID) ||
		(oPlayer != nil && oPlayer.ID == playerID) {
		return domain.JoinGameResult{}, &domain.AddGamePlayerError{
			Err:      domain.ErrAlreadyJoined,
			PlayerID: playerID,
			GameID:   gameID,
		}
	}

	side := domain.NoneSide
	switch {
	case xPlayer == nil:
		side = domain.XSide
	case oPlayer == nil:
		side = domain.OSide
	default:
		return domain.JoinGameResult{}, &domain.AddGamePlayerError{
			Err:      domain.ErrAllPlacesAlreadyTaken,
			PlayerID: playerID,
			GameID:   gameID,
		}
	}

	err = uc.gameRepo.AddGamePlayer(ctx, gameID, playerID, side)
	if err != nil {
		return domain.JoinGameResult{}, err
	}

	ready := false
	if xPlayer != nil || oPlayer != nil {
		ready = true
	}

	return domain.JoinGameResult{
		Side:         side,
		ReadyToStart: ready,
	}, nil
}

func (uc *GameUC) StartGame(ctx context.Context, gameID uuid.UUID) error {
	g, err := uc.gameRepo.GetGame(ctx, gameID)
	if err != nil {
		return err
	}

	if g.State == domain.Finished {
		return &domain.GameErrorWithID{
			Err: domain.ErrGameFinished,
			ID:  gameID,
		}
	}

	if g.State == domain.Started {
		return &domain.GameErrorWithID{
			Err: domain.ErrGameAlreadyStarted,
			ID:  gameID,
		}
	}

	if g.XPlayer == nil || g.OPlayer == nil {
		return &domain.GameErrorWithID{
			Err: domain.ErrNotEnoughPlayers,
			ID:  gameID,
		}
	}

	err = uc.gameRepo.UpdateGameState(ctx, g.ID, domain.Started)
	if err != nil {
		return err
	}

	return nil
}

func (uc *GameUC) GetSide(ctx context.Context, gameID uuid.UUID, playerID domain.PlayerID) (domain.Side, error) {
	xPlayer, oPlayer, err := uc.gameRepo.GetPlayers(ctx, gameID)
	if err != nil {
		return domain.NoneSide, err
	}

	side := domain.NoneSide

	switch {
	case xPlayer != nil && xPlayer.ID == playerID:
		side = domain.XSide
	case oPlayer != nil && oPlayer.ID == playerID:
		side = domain.OSide
	default:
		if xPlayer != nil && oPlayer != nil {
			return domain.NoneSide, &domain.AddGamePlayerError{
				Err:      domain.ErrAllPlacesAlreadyTaken,
				PlayerID: playerID,
				GameID:   gameID,
			}
		}
	}

	return side, nil
}

/*func (uc *GameUC) ReadyGame(ctx context.Context, gameID uuid.UUID, playerID domain.PlayerID) error {

}*/

func (uc *GameUC) MakeMove(ctx context.Context, gameID uuid.UUID, move domain.Move) (domain.MakeMoveResult, error) {
	g, err := uc.gameRepo.GetGame(ctx, gameID)
	if err != nil {
		return domain.MakeMoveResult{}, err
	}

	if g.State == domain.Created {
		return domain.MakeMoveResult{}, &domain.GameErrorWithID{Err: domain.ErrGameNotStarted, ID: g.ID}
	}
	if g.State == domain.Finished {
		return domain.MakeMoveResult{}, &domain.GameErrorWithID{Err: domain.ErrGameFinished, ID: g.ID}
	}

	err = uc.gameMode.IterateGame(ctx, g, move)
	if err != nil {
		return domain.MakeMoveResult{}, err
	}

	err = uc.gameRepo.UpdateGame(ctx, g)
	if err != nil {
		return domain.MakeMoveResult{}, err
	}

	if g.State == domain.Finished {
		return domain.MakeMoveResult{GameFinished: true}, nil
	}

	return domain.MakeMoveResult{}, nil
}
