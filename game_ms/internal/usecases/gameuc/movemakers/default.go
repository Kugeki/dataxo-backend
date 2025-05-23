package movemakers

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"dataxo-backend-game-ms/internal/usecases/gameuc"
	"dataxo-backend-game-ms/pkg/slogdiscard"
	"log/slog"
)

type Default struct {
	Cfg domain.DisappearingModeConfig

	log *slog.Logger
}

func NewDefault(cfg domain.DisappearingModeConfig, log *slog.Logger) *Default {
	if log == nil {
		log = slogdiscard.Logger()
	}

	return &Default{Cfg: cfg, log: log}
}

func (r *Default) MakeMoveOnBoard(ctx context.Context, g *domain.Game, board gameuc.Board, move domain.Move) []domain.MoveEvent {
	events := make([]domain.MoveEvent, 0, 2)

	move.TimesUsed = board.GetMove(move.XY()).TimesUsed + 1

	events = append(events, domain.MoveEvent{
		Type: domain.PlaceMove,
		Move: move,
	})
	board.SetMove(move)
	g.Moves = append(g.Moves, move)

	limit := r.Cfg.PlayerFiguresLimit * 2

	if limit <= 0 {
		return events
	}

	limitedIndex := len(g.Moves) - limit - 1

	if limitedIndex < 0 {
		return events
	}

	events = append(events, domain.MoveEvent{
		Type: domain.RemoveMove,
		Move: g.Moves[limitedIndex],
	})
	g.Moves[limitedIndex].Side = domain.NoneSide
	board.SetMove(g.Moves[limitedIndex])

	return events
}
