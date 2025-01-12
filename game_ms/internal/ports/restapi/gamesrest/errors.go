package gamesrest

import (
	"errors"
	"fmt"
)

var (
	ErrWrongMessageType = errors.New("wrong message type")

	ErrGameNotStarted     = errors.New("game not started")
	ErrGameAlreadyStarted = errors.New("game already started")

	ErrInvalidPresenceAction = errors.New("invalid presence action")

	ErrInvalidReadinessAction = errors.New("invalid readiness action")
)

type PresenceActionError struct {
	Err    error
	Action string
}

func (e *PresenceActionError) Error() string {
	return fmt.Sprintf("action(%v): %v", e.Action, e.Err)
}

func (e *PresenceActionError) Unwrap() error {
	return e.Err
}

type ReadinessError struct {
	Err    error
	Action string
}

func (e *ReadinessError) Error() string {
	return fmt.Sprintf("action(%v): %v", e.Action, e.Err)
}

func (e *ReadinessError) Unwrap() error {
	return e.Err
}
