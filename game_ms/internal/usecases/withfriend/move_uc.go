package withfriend

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"dataxo-backend-game-ms/pkg/slogdiscard"
	"errors"
	"log/slog"
)

type MoveValidator interface {
	ValidateMove(ctx context.Context, game *domain.Game, move domain.Move) error
}

type WinChecker interface {
	CheckWin(ctx context.Context, game *domain.Game, boardSize BoardSize, move domain.Move) (domain.WinResult, error)
}

type DefaultMoveMaker struct {
	Cfg Config

	validator  MoveValidator
	winChecker WinChecker

	log *slog.Logger
}

func New(cfg Config, log *slog.Logger) (*DefaultMoveMaker, error) {
	err := ValidateConfig(cfg)
	if err != nil {
		return nil, err
	}

	if log == nil {
		log = slogdiscard.Logger()
	}

	validator := NewDefaultMoveValidator(cfg, log)
	winChecker := NewDefaultWinChecker(cfg.WinLineLength)

	return &DefaultMoveMaker{Cfg: cfg, validator: validator, winChecker: winChecker, log: log}, nil
}

func (r *DefaultMoveMaker) MakeMove(ctx context.Context, g *domain.Game, move domain.Move) error {
	if g == nil {
		return errors.New("game is nil")
	}

	err := r.validator.ValidateMove(ctx, g, move)
	if err != nil {
		return &domain.GameErrorWithID{Err: err, ID: g.ID}
	}

	r.MakeMoveOnBoard(ctx, g, move)

	boardSize := BoardSize{
		Width:  r.Cfg.BoardWidth,
		Height: r.Cfg.BoardHeight,
	}

	// todo: add win sequence
	checker := NewDefaultWinChecker(r.Cfg.WinLineLength)
	winResult, err := checker.CheckWin(ctx, g, boardSize, move)
	if err != nil {
		return &domain.GameErrorWithID{Err: err, ID: g.ID}
	}

	switch winResult.Side {
	case domain.NoneSide:
		return nil
	case domain.XSide:
		fallthrough
	case domain.OSide:
		g.Winner = winResult.Side
		g.State = domain.Finished
		// g.WinSequence = ...
		return nil
	default:
		err = &domain.MoveError{Err: domain.ErrInvalidSide, Move: move}
		return &domain.GameErrorWithID{Err: err, ID: g.ID}
	}
}

func (r *DefaultMoveMaker) MakeMoveOnBoard(ctx context.Context, g *domain.Game, move domain.Move) {
	g.Board[move.Y][move.X] = move.Side
	g.Moves = append(g.Moves, move)

	if r.Cfg.PlayerFiguresLimit == 0 {
		return
	}

	limit := r.Cfg.PlayerFiguresLimit
	limitedIndex := len(g.Moves) - limit
	limitedMove := g.Moves[limitedIndex]

	g.Board[limitedMove.Y][limitedMove.X] = domain.NoneSide
}

/*func (r *DefaultMoveMaker) MakeMoveWithoutValidate(ctx context.Context, g *domain.Game, move domain.Move) error {
	if g == nil {
		return errors.New("game is nil")
	}
}*/
