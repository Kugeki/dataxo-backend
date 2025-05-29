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

func (r *GameRepoMap) CreateGame(ctx context.Context, plID domain.PlayerID, side domain.Side,
	mode string, cfg domain.DisappearingModeConfig) (*domain.Game, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	g := &domain.Game{
		ID:          id,
		Mode:        mode,
		Config:      cfg,
		State:       domain.Created,
		Moves:       make([]domain.Move, 0),
		WinSequence: make([]domain.Move, 0),
	}

	player := domain.Player{ID: plID, Ready: false}

	switch side {
	case domain.XSide:
		g.XPlayer = &player
	case domain.OSide:
		g.OPlayer = &player
	default:
		return nil, domain.ErrInvalidSide
	}

	r.m[g.ID] = g

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

func (r *GameRepoMap) UpdateGameState(ctx context.Context, gameID uuid.UUID, state domain.State) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.m[gameID]
	if !ok {
		return &domain.GameErrorWithID{
			Err: domain.ErrNotFound,
			ID:  gameID,
		}
	}

	r.m[gameID].State = state
	return nil
}

func (r *GameRepoMap) UpdateGame(ctx context.Context, g *domain.Game) error {
	// todo: make changes atomic. i.e. make change categories enum
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.m[g.ID]
	if !ok {
		return &domain.GameErrorWithID{
			Err: domain.ErrNotFound,
			ID:  g.ID,
		}
	}

	r.m[g.ID] = g
	return nil
}

func (r *GameRepoMap) GetPlayers(ctx context.Context, gameID uuid.UUID) (x *domain.Player, o *domain.Player, err error) {
	g, err := r.GetGame(ctx, gameID)
	if err != nil {
		return nil, nil, err
	}

	return g.XPlayer, g.OPlayer, nil
}

func (r *GameRepoMap) AddGamePlayer(ctx context.Context, gameID uuid.UUID, playerID domain.PlayerID, side domain.Side) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	g, ok := r.m[gameID]
	if !ok || g == nil {
		return &domain.GameErrorWithID{
			Err: domain.ErrNotFound,
			ID:  gameID,
		}
	}

	player := domain.Player{ID: playerID, Ready: false}
	switch side {
	case domain.XSide:
		g.XPlayer = &player
	case domain.OSide:
		g.OPlayer = &player
	default:
		return &domain.AddGamePlayerError{
			Err:      domain.ErrInvalidSide,
			PlayerID: playerID,
			GameID:   gameID,
		}
	}

	r.m[gameID] = g

	return nil
}
