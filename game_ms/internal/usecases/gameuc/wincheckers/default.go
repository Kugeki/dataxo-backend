package wincheckers

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"dataxo-backend-game-ms/internal/usecases/gameuc"
	"errors"
	"fmt"
)

func NoWinner() domain.WinResult {
	return domain.NoWinner()
}

type Default struct {
	Game          *domain.Game
	Board         gameuc.Board
	BoardSize     gameuc.BoardSize
	Move          domain.Move
	WinLineLength int

	winSequenceBuffer []domain.Move
}

func NewDefault(winLineLength int) *Default {
	return &Default{WinLineLength: winLineLength, winSequenceBuffer: make([]domain.Move, winLineLength)}
}

type Check func(ctx context.Context) domain.WinResult

func (c *Default) CheckWin(ctx context.Context, game *domain.Game, board gameuc.Board,
	boardSize gameuc.BoardSize, move domain.Move) (domain.WinResult, error) {

	if game == nil {
		return NoWinner(), errors.New("game is nil")
	}

	c.Game = game
	c.Board = board
	c.BoardSize = boardSize
	c.Move = move

	checks := []Check{
		c.CheckWinUpDown,
		c.CheckWinLeftRight,
		c.CheckWinLeftDown,
		c.CheckWinRightDown,
	}

	var winSides []domain.WinSide
	var winSequence []domain.Move

	for _, check := range checks {
		winResult := check(ctx)

		if winResult.Side != domain.NoneWin {
			winSides = append(winSides, winResult.Side)
			winSequence = winResult.Sequence
		}
	}

	switch len(winSides) {
	case 0:
		break
	case 1:
		return domain.WinResult{Side: winSides[0], Sequence: winSequence}, nil
	default:
		return NoWinner(), fmt.Errorf("wrong count(%v) of win sides(%v)", len(winSides), winSides)
	}

	h, w := boardSize.Height, boardSize.Width
	if h*w == min(len(game.Moves)+1, game.Config.PlayerFiguresLimit) {
		return domain.WinResult{
			Side:     domain.Draw,
			Sequence: make([]domain.Move, 0),
		}, nil
	}

	return NoWinner(), nil
}

func (c *Default) CheckWinUpDown(ctx context.Context) domain.WinResult {
	x, y := c.Move.XY()
	side := c.Move.Side

	for ; y > 0; y-- {
		if c.Board.GetMove(x, y-1).Side != side {
			break
		}
	}

	toY := y + c.WinLineLength - 1
	if toY >= c.BoardSize.Height {
		return NoWinner()
	}

	sequence := c.GetSequence()
	for ; y <= toY; y++ {
		move := c.Board.GetMove(x, y)
		if move.Side != side {
			return NoWinner()
		}
		sequence = append(sequence, move)
	}

	return domain.WinResult{
		Side:     side.ToWinSide(),
		Sequence: sequence,
	}
}

func (c *Default) CheckWinLeftRight(ctx context.Context) domain.WinResult {
	x, y := c.Move.XY()
	side := c.Move.Side

	for ; x > 0; x-- {
		if c.Board.GetMove(x-1, y).Side != side {
			break
		}
	}

	toX := x + c.WinLineLength - 1
	if toX >= c.BoardSize.Width {
		return NoWinner()
	}

	sequence := c.GetSequence()
	for ; x <= toX; x++ {
		move := c.Board.GetMove(x, y)
		if move.Side != side {
			return NoWinner()
		}
		sequence = append(sequence, move)
	}

	return domain.WinResult{
		Side:     side.ToWinSide(),
		Sequence: sequence,
	}
}

func (c *Default) CheckWinLeftDown(ctx context.Context) domain.WinResult {
	x, y := c.Move.XY()
	side := c.Move.Side

	for y > 0 && x > 0 {
		if c.Board.GetMove(x-1, y-1).Side != side {
			break
		}

		y--
		x--
	}

	toY := y + c.WinLineLength - 1
	if toY >= c.BoardSize.Height {
		return NoWinner()
	}

	toX := x + c.WinLineLength - 1
	if toX >= c.BoardSize.Width {
		return NoWinner()
	}

	sequence := c.GetSequence()
	for y <= toY && x <= toX {
		move := c.Board.GetMove(x, y)
		if move.Side != side {
			return NoWinner()
		}
		sequence = append(sequence, move)

		y++
		x++
	}
	return domain.WinResult{
		Side:     side.ToWinSide(),
		Sequence: sequence,
	}
}

func (c *Default) CheckWinRightDown(ctx context.Context) domain.WinResult {
	x, y := c.Move.XY()
	side := c.Move.Side

	for y > 0 && x < c.BoardSize.Width-1 {
		if c.Board.GetMove(x+1, y-1).Side != side {
			break
		}

		y--
		x++
	}

	toY := y + c.WinLineLength - 1
	if toY >= c.BoardSize.Height {
		return NoWinner()
	}

	toX := x - c.WinLineLength + 1
	if toX >= c.BoardSize.Width {
		return NoWinner()
	}

	sequence := c.GetSequence()
	for y <= toY && x >= toX {
		move := c.Board.GetMove(x, y)
		if move.Side != side {
			return NoWinner()
		}
		sequence = append(sequence, move)

		y++
		x--
	}
	return domain.WinResult{
		Side:     side.ToWinSide(),
		Sequence: sequence,
	}
}

func (c *Default) GetSequence() []domain.Move {
	return make([]domain.Move, 0, c.WinLineLength)
}
