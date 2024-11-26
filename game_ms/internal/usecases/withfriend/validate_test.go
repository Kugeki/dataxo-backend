package withfriend

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultMoveValidator_ValidateMoveCoords(t *testing.T) {
	tcases := []struct {
		Name      string
		BoardSize BoardSize
		MoveX     int
		MoveY     int
		Err       error
	}{
		{
			Name: "valid move",
			BoardSize: BoardSize{
				Width:  10,
				Height: 5,
			},
			MoveX: 3,
			MoveY: 2,
			Err:   nil,
		},
		{
			Name: "overflow X",
			BoardSize: BoardSize{
				Width:  10,
				Height: 5,
			},
			MoveX: 13,
			MoveY: 2,
			Err:   domain.ErrMoveOutOfBoard,
		},
		{
			Name: "overflow Y",
			BoardSize: BoardSize{
				Width:  10,
				Height: 5,
			},
			MoveX: 3,
			MoveY: 80,
			Err:   domain.ErrMoveOutOfBoard,
		},
		{
			Name: "negative X",
			BoardSize: BoardSize{
				Width:  10,
				Height: 5,
			},
			MoveX: -10,
			MoveY: 3,
			Err:   domain.ErrMoveOutOfBoard,
		},
		{
			Name: "negative Y",
			BoardSize: BoardSize{
				Width:  10,
				Height: 5,
			},
			MoveX: 5,
			MoveY: -1,
			Err:   domain.ErrMoveOutOfBoard,
		},
		{
			Name: "zeroed X",
			BoardSize: BoardSize{
				Width:  10,
				Height: 5,
			},
			MoveX: 0,
			MoveY: 3,
			Err:   domain.ErrMoveOutOfBoard,
		},
		{
			Name: "zeroed Y",
			BoardSize: BoardSize{
				Width:  10,
				Height: 5,
			},
			MoveX: 5,
			MoveY: 0,
			Err:   domain.ErrMoveOutOfBoard,
		},
	}

	for _, tc := range tcases {
		validator := NewDefaultMoveValidator(Config{
			PlayerFiguresLimit: 4,
			WinLineLength:      3,
			BoardWidth:         tc.BoardSize.Width,
			BoardHeight:        tc.BoardSize.Width,
		}, nil)

		err := validator.ValidateMoveCoords(context.Background(), tc.BoardSize, tc.MoveX, tc.MoveY)
		assert.ErrorIs(t, err, tc.Err)
	}
}
