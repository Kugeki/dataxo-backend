package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type Game struct {
	ID          uuid.UUID
	Mode        string
	Config      DisappearingModeConfig
	State       State
	Moves       []Move
	XPlayer     *Player
	OPlayer     *Player
	WinSequence []Move
	Winner      WinSide
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
	ID        int  `json:"id"`
	InGameID  int  `json:"in_game_id"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	TimesUsed int  `json:"times_used"`
	Side      Side `json:"side"`
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
	return fmt.Sprintf("move_id(%v) ingame_id(%v) max_ingame_id(%v): %v",
		e.Move.ID, e.Move.InGameID, e.MaxInGameID, e.Err)
}

func (e *MoveErrorWithInGameID) Unwrap() error {
	return e.Err
}

type PlayerID struct {
	ClientID string
}

type Player struct {
	ID    PlayerID
	Ready bool
}

type PlayerError struct {
	Err      error
	PlayerID PlayerID
}

func (e *PlayerError) Error() string {
	return fmt.Sprintf("player_id(%v): %v",
		e.PlayerID, e.Err)
}

func (e *PlayerError) Unwrap() error {
	return e.Err
}

type AddGamePlayerError struct {
	Err      error
	PlayerID PlayerID
	Side     Side
	GameID   uuid.UUID
}

func (e *AddGamePlayerError) Error() string {
	return fmt.Sprintf("player_id(%v) side(%v) game_id(%v): %v",
		e.PlayerID, e.Side, e.GameID, e.Err)
}

func (e *AddGamePlayerError) Unwrap() error {
	return e.Err
}

type ModeParams struct {
	MySide SideRequest
}

type SideRequest int

const (
	RandomSideRequest SideRequest = iota
	XSideRequest
	OSideRequest
)

type Side int

const (
	NoneSide Side = iota
	XSide
	OSide
)

func (s Side) ToWinSide() WinSide {
	switch s {
	case NoneSide:
		return NoneWin
	case XSide:
		return XWin
	case OSide:
		return OWin
	default:
		// todo: process
		return NoneWin
	}
}

func NoneSideMove() Move {
	return Move{Side: NoneSide}
}

type WinResult struct {
	Side     WinSide
	Sequence []Move
}

type WinSide int

const (
	NoneWin WinSide = iota
	XWin
	OWin
	Draw
)

func (r WinResult) IsNoWinner() bool {
	return r.Side == NoneWin
}

func NoWinner() WinResult {
	return WinResult{Side: NoneWin}
}

type JoinGameResult struct {
	Side         Side
	ReadyToStart bool
}

type MoveEventType int

const (
	PlaceMove MoveEventType = iota
	RemoveMove
	HeatMove
	BlockMove
)

type MoveEvent struct {
	Type MoveEventType
	Move Move
}

type MakeMoveResult struct {
	GameFinished bool
	Events       []MoveEvent
}
