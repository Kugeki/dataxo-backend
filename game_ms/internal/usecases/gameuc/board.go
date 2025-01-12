package gameuc

import "dataxo-backend-game-ms/internal/domain"

type Pos struct {
	X, Y int
}

func MovePos(move domain.Move) Pos {
	return Pos{X: move.X, Y: move.Y}
}

type Board map[Pos]domain.Move

func (b Board) SetMove(move domain.Move) {
	b[MovePos(move)] = move
}

func (b Board) GetMove(x, y int) domain.Move {
	res, _ := b[Pos{X: x, Y: y}]
	return res
}

func NewBoard(moves []domain.Move, figuresLimit int) Board {
	if figuresLimit == 0 {
		figuresLimit = len(moves)
	}

	board := Board(make(map[Pos]domain.Move, len(moves)))

	for i := max(len(moves)-figuresLimit*2, 0); i < len(moves); i++ {
		board.SetMove(moves[i])
	}

	return board
}

type BoardSize struct {
	Width  int
	Height int
}
