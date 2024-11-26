package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type Game struct {
	ID          uuid.UUID
	Mode        string
	State       State
	Moves       []Move
	Board       [][]Side
	XPlayer     Player
	OPlayer     Player
	WinSequence []Move
	Winner      Side
}

type GameErrorWithID struct {
	Err error
	ID  uuid.UUID
}

func (e *GameErrorWithID) Error() string {
	return fmt.Sprintf("game with id %v: %v", e.ID, e.Err)
}

func (e *GameErrorWithID) Unwrap() error {
	return e.Err
}

const ModeWithFriend = "with-friend"

type State int

const (
	Created State = iota
	Started
	Finished
)

func (s State) String() string {
	switch s {
	case Created:
		return "created"
	case Started:
		return "started"
	case Finished:
		return "finished"
	default:
		return "invalid"
	}
}

type Move struct {
	ID       int
	InGameID int
	X        int
	Y        int
	Side     Side
}

func (m Move) XY() (int, int) {
	return m.X, m.Y
}

type MoveError struct {
	Err  error
	Move Move
}

func (e *MoveError) Error() string {
	return fmt.Sprintf("move id(%v) x(%v) y(%v): %v",
		e.Move.ID, e.Move.X, e.Move.Y, e.Err)
}

func (e *MoveError) Unwrap() error {
	return e.Err
}

type MoveErrorWithInGameID struct {
	Err         error
	Move        Move
	MaxInGameID int
}

func (e *MoveErrorWithInGameID) Error() string {
	return fmt.Sprintf("move id(%v) ingame_id(%v) max_ingame_id(%v): %v",
		e.Move.ID, e.Move.InGameID, e.MaxInGameID, e.Err)
}

func (e *MoveErrorWithInGameID) Unwrap() error {
	return e.Err
}

type Player struct {
	RemoteAddr string
}

type ModeParams struct {
	MySide Side
}

type Side int

const (
	NoneSide Side = iota - 1
	XSide
	OSide
)

type WinResult struct {
	Side     Side
	Sequence []Move
}

func (r WinResult) IsNoWinner() bool {
	return r.Side == NoneSide
}

func NoWinner() WinResult {
	return WinResult{Side: NoneSide}
}
