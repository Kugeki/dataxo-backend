package domain

import "errors"

var (
	ErrNotFound = errors.New("not found")

	ErrInvalidSide    = errors.New("invalid side")
	ErrInvalidWinSide = errors.New("invalid win side")

	ErrGameNotStarted     = errors.New("game not started")
	ErrGameAlreadyStarted = errors.New("game already started")
	ErrGameFinished       = errors.New("game finished")

	ErrGameIsNil = errors.New("game is nil")

	ErrPlaceAlreadyTaken   = errors.New("place already taken")
	ErrMoveOutOfBoard      = errors.New("move is out of board")
	ErrInvalidMoveInGameID = errors.New("invalid move ingame id")
	ErrInvalidSideTurn     = errors.New("now is not your side turn")

	ErrAlreadyJoined         = errors.New("already joined")
	ErrAllPlacesAlreadyTaken = errors.New("all places already taken in this game")
	ErrNotEnoughPlayers      = errors.New("not enough players")

	ErrNegativePlayerFiguresLimit    = errors.New("player figures limit is negative")
	ErrNegativeOrZeroedWinLineLength = errors.New("win line length is negative or equals to zero")
	ErrNegativeOrZeroedBoardWidth    = errors.New("board width is negative or equals to zero")
	ErrNegativeOrZeroedBoardHeight   = errors.New("board height is negative or equals to zero")
)
