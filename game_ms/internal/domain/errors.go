package domain

import "errors"

var (
	ErrNotFound = errors.New("not found")

	ErrInvalidSide = errors.New("invalid side")

	ErrPlaceAlreadyTaken   = errors.New("place already taken")
	ErrMoveOutOfBoard      = errors.New("move is out of board")
	ErrInvalidBoardSize    = errors.New("invalid board size")
	ErrInvalidMoveInGameID = errors.New("invalid move ingame id")
)
