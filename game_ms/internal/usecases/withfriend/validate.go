package withfriend

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"dataxo-backend-game-ms/pkg/slogdiscard"
	"log/slog"
)

type DefaultMoveValidator struct {
	log *slog.Logger

	Cfg Config
}

func NewDefaultMoveValidator(cfg Config, log *slog.Logger) *DefaultMoveValidator {
	if log == nil {
		log = slogdiscard.Logger()
	}

	return &DefaultMoveValidator{Cfg: cfg, log: log}
}

func (v *DefaultMoveValidator) ValidateMove(ctx context.Context, game *domain.Game, move domain.Move) error {
	if game.Board[move.Y][move.X] != domain.NoneSide {
		return &domain.MoveError{
			Err:  domain.ErrPlaceAlreadyTaken,
			Move: move,
		}
	}

	maxInGameID := v.GetMaxMoveInGameID(ctx, game.Moves)

	err := v.ValidateMoveInGameID(ctx, move.InGameID, maxInGameID, game.Moves)
	if err != nil {
		return &domain.MoveErrorWithInGameID{Err: err, Move: move, MaxInGameID: maxInGameID}
	}

	boardSize, err := v.GetBoardSize(ctx, game.Board)
	if err != nil {
		return err
	}

	err = v.ValidateBoardSize(ctx, boardSize)
	if err != nil {
		return err
	}

	err = v.ValidateMoveCoords(ctx, boardSize, move.X, move.Y)
	if err != nil {
		return &domain.MoveError{Err: err, Move: move}
	}

	err = v.ValidateSide(ctx, move.Side)
	if err != nil {
		return &domain.MoveError{Err: err, Move: move}
	}

	return nil
}

type BoardSize struct {
	Width  int
	Height int
}

func (v *DefaultMoveValidator) GetMaxMoveInGameID(ctx context.Context, moves []domain.Move) int {
	if len(moves) <= 0 {
		return 0
	}

	maxID := moves[0].InGameID
	for _, v := range moves {
		maxID = max(maxID, v.InGameID)
	}

	return maxID
}

func (v *DefaultMoveValidator) ValidateMoveInGameID(ctx context.Context, id, maxID int, moves []domain.Move) error {
	if len(moves) <= 0 {
		if id != 0 {
			return domain.ErrInvalidMoveInGameID
		}

		return nil
	}

	if id != maxID+1 {
		return domain.ErrInvalidMoveInGameID
	}

	if id != len(moves) {
		return domain.ErrInvalidMoveInGameID
	}

	return nil
}

func (v *DefaultMoveValidator) GetBoardSize(ctx context.Context, board [][]domain.Side) (BoardSize, error) {
	height := len(board)

	if height <= 0 {
		return BoardSize{}, domain.ErrInvalidBoardSize
	}

	width := len(board[0])
	for i := range board {
		currWidth := len(board[i])

		if currWidth != width {
			v.log.Error(
				"some of the board width don't equal to entire width",
				slog.Int("y", i),
				slog.Int("entire width", width),
				slog.Int("current width", currWidth),
			)
		}

		width = min(width, currWidth)
	}

	return BoardSize{
		Width:  width,
		Height: height,
	}, nil
}

func (v *DefaultMoveValidator) ValidateBoardSize(ctx context.Context, size BoardSize) error {
	if size.Height != v.Cfg.BoardHeight || size.Width != v.Cfg.BoardWidth ||
		size.Height <= 0 || size.Width <= 0 {
		return domain.ErrInvalidBoardSize
	}

	return nil
}

func (v *DefaultMoveValidator) ValidateMoveCoords(ctx context.Context, boardSize BoardSize, x, y int) error {
	if y >= boardSize.Height || y <= 0 ||
		x >= boardSize.Width || x <= 0 {
		return domain.ErrMoveOutOfBoard
	}
	return nil
}

func (v *DefaultMoveValidator) ValidateSide(ctx context.Context, side domain.Side) error {
	switch side {
	case domain.NoneSide:
	case domain.XSide:
	case domain.OSide:
	default:
		return domain.ErrInvalidSide
	}

	return nil
}
