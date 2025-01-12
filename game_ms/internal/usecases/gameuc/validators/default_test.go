package validators

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"dataxo-backend-game-ms/internal/usecases/gameuc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultMoveValidator_ValidateMoveCoords(t *testing.T) {
	tcases := []struct {
		Name      string
		BoardSize gameuc.BoardSize
		MoveX     int
		MoveY     int
		Err       error
	}{
		{
			Name: "valid move",
			BoardSize: gameuc.BoardSize{
				Width:  10,
				Height: 5,
			},
			MoveX: 3,
			MoveY: 2,
			Err:   nil,
		},
		{
			Name: "overflow X",
			BoardSize: gameuc.BoardSize{
				Width:  10,
				Height: 5,
			},
			MoveX: 13,
			MoveY: 2,
			Err:   domain.ErrMoveOutOfBoard,
		},
		{
			Name: "overflow Y",
			BoardSize: gameuc.BoardSize{
				Width:  10,
				Height: 5,
			},
			MoveX: 3,
			MoveY: 80,
			Err:   domain.ErrMoveOutOfBoard,
		},
		{
			Name: "negative X",
			BoardSize: gameuc.BoardSize{
				Width:  10,
				Height: 5,
			},
			MoveX: -10,
			MoveY: 3,
			Err:   domain.ErrMoveOutOfBoard,
		},
		{
			Name: "negative Y",
			BoardSize: gameuc.BoardSize{
				Width:  10,
				Height: 5,
			},
			MoveX: 5,
			MoveY: -1,
			Err:   domain.ErrMoveOutOfBoard,
		},
	}

	for _, tc := range tcases {
		validator := NewDefault(domain.DisappearingModeConfig{
			PlayerFiguresLimit: 4,
			WinLineLength:      3,
			BoardWidth:         tc.BoardSize.Width,
			BoardHeight:        tc.BoardSize.Width,
		}, nil)

		err := validator.ValidateMoveCoords(context.Background(), tc.BoardSize, tc.MoveX, tc.MoveY)
		assert.ErrorIs(t, err, tc.Err)
	}
}
