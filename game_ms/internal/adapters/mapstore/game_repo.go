package mapstore

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"github.com/google/uuid"
	"sync"
)

type GameRepoMap struct {
	m  map[uuid.UUID]*domain.Game
	mu sync.Mutex
}

func NewGameRepo() *GameRepoMap {
	return &GameRepoMap{m: make(map[uuid.UUID]*domain.Game), mu: sync.Mutex{}}
}

func (r *GameRepoMap) CreateGame(ctx context.Context, pl domain.Player,
	mode string, params domain.ModeParams) (*domain.Game, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	g := &domain.Game{
		ID:    id,
		Mode:  mode,
		State: domain.Created,
		Moves: make([]domain.Move, 0),
	}

	switch params.MySide {
	case domain.XSide:
		g.XPlayer = pl
	case domain.OSide:
		g.OPlayer = pl
	default:
		return nil, domain.ErrInvalidSide
	}

	return g, nil
}

func (r *GameRepoMap) GetGame(ctx context.Context, gameID uuid.UUID) (*domain.Game, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	g, ok := r.m[gameID]
	if !ok || g == nil {
		return nil, &domain.GameErrorWithID{
			Err: domain.ErrNotFound,
			ID:  gameID,
		}
	}

	return g, nil
}
