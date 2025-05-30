package modes

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"dataxo-backend-game-ms/internal/usecases/gameuc"
	"dataxo-backend-game-ms/internal/usecases/gameuc/movemakers"
	"dataxo-backend-game-ms/internal/usecases/gameuc/validators"
	"dataxo-backend-game-ms/internal/usecases/gameuc/wincheckers"
	"dataxo-backend-game-ms/pkg/slogdiscard"
	"log/slog"
)

type MoveValidator interface {
	ValidateMove(ctx context.Context, game *domain.Game, board gameuc.Board, move domain.Move) error
}

type WinChecker interface {
	CheckWin(ctx context.Context, game *domain.Game, board gameuc.Board, boardSize gameuc.BoardSize, move domain.Move) (domain.WinResult, error)
}

type MoveMaker interface {
	MakeMoveOnBoard(ctx context.Context, g *domain.Game, board gameuc.Board, move domain.Move) []domain.MoveEvent
}

type DisappearingMode struct {
	Cfg domain.DisappearingModeConfig

	validator MoveValidator
	checker   WinChecker
	moveMaker MoveMaker

	log *slog.Logger
}

func NewDisappearingMode(cfg domain.DisappearingModeConfig, log *slog.Logger) (*DisappearingMode, error) {
	err := domain.ValidateDisappearingModeConfig(cfg)
	if err != nil {
		return nil, err
	}

	validator := validators.NewDefault(cfg, log)
	winChecker := wincheckers.NewDefault(cfg.WinLineLength)
	moveMaker := movemakers.NewDefault(cfg, log)
	if err != nil {
		return nil, err
	}

	return &DisappearingMode{Cfg: cfg, validator: validator, checker: winChecker,
		moveMaker: moveMaker, log: slogdiscard.LoggerIfNil(log)}, nil
}

func (m *DisappearingMode) IterateGame(ctx context.Context, g *domain.Game, move domain.Move) (domain.MakeMoveResult, error) {
	if g == nil {
		return domain.MakeMoveResult{}, domain.ErrGameIsNil
	}

	if g.State == domain.Created {
		return domain.MakeMoveResult{}, &domain.GameErrorWithID{Err: domain.ErrGameNotStarted, ID: g.ID}
	}
	if g.State == domain.Finished {
		return domain.MakeMoveResult{}, &domain.GameErrorWithID{Err: domain.ErrGameFinished, ID: g.ID}
	}

	board := gameuc.NewBoard(g.Moves, m.Cfg.PlayerFiguresLimit)

	err := m.validator.ValidateMove(ctx, g, board, move)
	if err != nil {
		return domain.MakeMoveResult{}, &domain.GameErrorWithID{Err: err, ID: g.ID}
	}

	moveEvents := m.moveMaker.MakeMoveOnBoard(ctx, g, board, move)

	boardSize := gameuc.BoardSize{
		Width:  m.Cfg.BoardWidth,
		Height: m.Cfg.BoardHeight,
	}

	// todo: add win sequence
	winResult, err := m.checker.CheckWin(ctx, g, board, boardSize, move)
	if err != nil {
		return domain.MakeMoveResult{}, &domain.GameErrorWithID{Err: err, ID: g.ID}
	}

	switch winResult.Side {
	case domain.NoneWin:
		return domain.MakeMoveResult{
			GameFinished: false,
			Events:       moveEvents,
		}, nil
	case domain.XWin, domain.OWin, domain.Draw:
		g.Winner = winResult.Side
		g.State = domain.Finished
		g.WinSequence = winResult.Sequence
		return domain.MakeMoveResult{
			GameFinished: true,
			Events:       moveEvents,
		}, nil
	default:
		err = &domain.MoveError{Err: domain.ErrInvalidWinSide, Move: move}
		return domain.MakeMoveResult{}, &domain.GameErrorWithID{Err: err, ID: g.ID}
	}
}

func (m *DisappearingMode) GetConfig() domain.DisappearingModeConfig {
	return m.Cfg
}
