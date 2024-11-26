package withfriend

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"errors"
	"fmt"
)

func NoWinner() domain.WinResult {
	return domain.NoWinner()
}

type DefaultWinChecker struct {
	Game          *domain.Game
	BoardSize     BoardSize
	Move          domain.Move
	WinLineLength int

	winSequenceBuffer []domain.Move
}

func NewDefaultWinChecker(winLineLength int) *DefaultWinChecker {
	return &DefaultWinChecker{WinLineLength: winLineLength, winSequenceBuffer: make([]domain.Move, winLineLength)}
}

type Check func(ctx context.Context) domain.WinResult

func (c *DefaultWinChecker) CheckWin(ctx context.Context, game *domain.Game,
	boardSize BoardSize, move domain.Move) (domain.WinResult, error) {

	if game == nil {
		return NoWinner(), errors.New("game is nil")
	}

	c.Game = game
	c.BoardSize = boardSize
	c.Move = move

	checks := []Check{
		c.CheckWinUp,
		c.CheckWinDown,
		c.CheckWinRight,
		c.CheckWinLeft,
		c.CheckWinUpRight,
		c.CheckWinUpLeft,
		c.CheckWinDownRight,
		c.CheckWinDownLeft,
	}

	// todo: add win sequence

	var winSides []domain.Side
	var winSequence []domain.Move

	for _, check := range checks {
		winResult := check(ctx)

		if winResult.Side != domain.NoneSide {
			winSides = append(winSides, winResult.Side)
			winSequence = winResult.Sequence
		}
	}

	switch len(winSides) {
	case 0:
		return NoWinner(), nil
	case 1:
		return domain.WinResult{Side: winSides[0], Sequence: winSequence}, nil
	default:
		return NoWinner(), fmt.Errorf("wrong count(%v) of win sides(%v)", len(winSides), winSides)
	}
}

func (c *DefaultWinChecker) CheckWinUp(ctx context.Context) domain.WinResult {
	x, y := c.Move.XY()
	side := c.Move.Side

	toY := y + c.WinLineLength - 1
	if toY >= c.BoardSize.Height {
		return NoWinner()
	}

	sequence := make([]domain.Move, 0, c.WinLineLength)
	for ; y <= toY; y++ {
		if c.Game.Board[y][x] != side {
			return NoWinner()
		}
		sequence = append(sequence)
	}

	return domain.WinResult{
		Side:     side,
		Sequence: nil,
	}
}

func (c *DefaultWinChecker) CheckWinDown(ctx context.Context) domain.WinResult {
	x, y := c.Move.XY()
	side := c.Move.Side

	toY := y - c.WinLineLength + 1
	if toY < 0 {
		return NoWinner()
	}

	for ; y >= toY; y-- {
		if c.Game.Board[y][x] != side {
			return NoWinner()
		}
	}

	return domain.WinResult{
		Side:     side,
		Sequence: nil,
	}
}

func (c *DefaultWinChecker) CheckWinRight(ctx context.Context) domain.WinResult {
	x, y := c.Move.XY()
	side := c.Move.Side

	toX := x + c.WinLineLength - 1
	if toX >= c.BoardSize.Width {
		return NoWinner()
	}

	for ; x <= toX; x++ {
		if c.Game.Board[y][x] != side {
			return NoWinner()
		}
	}

	return domain.WinResult{
		Side:     side,
		Sequence: nil,
	}
}

func (c *DefaultWinChecker) CheckWinLeft(ctx context.Context) domain.WinResult {
	x, y := c.Move.XY()
	side := c.Move.Side

	toX := x - c.WinLineLength + 1
	if toX < 0 {
		return NoWinner()
	}

	for ; x >= toX; x-- {
		if c.Game.Board[y][x] != side {
			return NoWinner()
		}
	}

	return domain.WinResult{
		Side:     side,
		Sequence: nil,
	}
}

func (c *DefaultWinChecker) CheckWinUpRight(ctx context.Context) domain.WinResult {
	x, y := c.Move.XY()
	side := c.Move.Side

	toY := y + c.WinLineLength - 1
	if toY >= c.BoardSize.Height {
		return NoWinner()
	}

	toX := x + c.WinLineLength - 1
	if toX >= c.BoardSize.Width {
		return NoWinner()
	}

	for y <= toY && x <= toX {
		if c.Game.Board[y][x] != side {
			return NoWinner()
		}
		y++
		x++
	}
	return domain.WinResult{
		Side:     side,
		Sequence: nil,
	}
}

func (c *DefaultWinChecker) CheckWinUpLeft(ctx context.Context) domain.WinResult {
	x, y := c.Move.XY()
	side := c.Move.Side

	toY := y + c.WinLineLength - 1
	if toY >= c.BoardSize.Height {
		return NoWinner()
	}

	toX := x - c.WinLineLength + 1
	if toX < 0 {
		return NoWinner()
	}

	for y <= toY && x >= toX {
		if c.Game.Board[y][x] != side {
			return NoWinner()
		}
		y++
		x--
	}
	return domain.WinResult{
		Side:     side,
		Sequence: nil,
	}
}

func (c *DefaultWinChecker) CheckWinDownRight(ctx context.Context) domain.WinResult {
	x, y := c.Move.XY()
	side := c.Move.Side

	toY := y - c.WinLineLength + 1
	if toY < 0 {
		return NoWinner()
	}

	toX := x + c.WinLineLength - 1
	if toX >= c.BoardSize.Width {
		return NoWinner()
	}

	for y >= toY && x <= toX {
		if c.Game.Board[y][x] != side {
			return NoWinner()
		}
		y--
		x++
	}
	return domain.WinResult{
		Side:     side,
		Sequence: nil,
	}
}

func (c *DefaultWinChecker) CheckWinDownLeft(ctx context.Context) domain.WinResult {
	x, y := c.Move.XY()
	side := c.Move.Side

	toY := y - c.WinLineLength + 1
	if toY < 0 {
		return NoWinner()
	}

	toX := x - c.WinLineLength + 1
	if toX < 0 {
		return NoWinner()
	}

	for y >= toY && x >= toX {
		if c.Game.Board[y][x] != side {
			return NoWinner()
		}
		y--
		x--
	}
	return domain.WinResult{
		Side:     side,
		Sequence: nil,
	}
}
