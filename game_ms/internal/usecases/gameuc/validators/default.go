package validators

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"dataxo-backend-game-ms/internal/usecases/gameuc"
	"dataxo-backend-game-ms/pkg/slogdiscard"
	"log/slog"
)

type Default struct {
	log *slog.Logger

	Cfg domain.DisappearingModeConfig
}

func NewDefault(cfg domain.DisappearingModeConfig, log *slog.Logger) *Default {
	if log == nil {
		log = slogdiscard.Logger()
	}

	return &Default{Cfg: cfg, log: log}
}

func (v *Default) ValidateMove(ctx context.Context, game *domain.Game, board gameuc.Board, move domain.Move) error {
	if board.GetMove(move.XY()).Side != domain.NoneSide {
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

	err = v.ValidateSideTurn(ctx, game.Moves, move.Side)
	if err != nil {
		return &domain.MoveError{Err: err, Move: move}
	}

	boardSize := v.GetBoardSize()
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

func (v *Default) GetMaxMoveInGameID(ctx context.Context, moves []domain.Move) int {
	if len(moves) <= 0 {
		return 0
	}

	maxID := moves[0].InGameID
	for _, v := range moves {
		maxID = max(maxID, v.InGameID)
	}

	return maxID
}

func (v *Default) ValidateMoveInGameID(ctx context.Context, id, maxID int, moves []domain.Move) error {
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

func (v *Default) ValidateSideTurn(ctx context.Context, moves []domain.Move, side domain.Side) error {
	if len(moves) < 1 {
		if side == domain.XSide {
			return nil
		}
		return domain.ErrInvalidSideTurn
	}

	prevSide := moves[len(moves)-1].Side
	if prevSide == side {
		return domain.ErrInvalidSideTurn
	}

	return nil
}

func (v *Default) GetBoardSize() gameuc.BoardSize {
	return gameuc.BoardSize{
		Width:  v.Cfg.BoardWidth,
		Height: v.Cfg.BoardHeight,
	}
}

func (v *Default) ValidateMoveCoords(ctx context.Context, boardSize gameuc.BoardSize, x, y int) error {
	if y >= boardSize.Height || y < 0 ||
		x >= boardSize.Width || x < 0 {
		return domain.ErrMoveOutOfBoard
	}
	return nil
}

func (v *Default) ValidateSide(ctx context.Context, side domain.Side) error {
	switch side {
	case domain.NoneSide:
	case domain.XSide:
	case domain.OSide:
	default:
		return domain.ErrInvalidSide
	}

	return nil
}
