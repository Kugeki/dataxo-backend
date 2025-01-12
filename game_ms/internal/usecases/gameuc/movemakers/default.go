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

func (r *Default) MakeMoveOnBoard(ctx context.Context, g *domain.Game, board gameuc.Board, move domain.Move) {
	board.SetMove(move)
	g.Moves = append(g.Moves, move)

	limit := r.Cfg.PlayerFiguresLimit * 2

	if limit <= 0 {
		return
	}

	limitedIndex := len(g.Moves) - limit - 1

	if limitedIndex < 0 {
		return
	}

	g.Moves[limitedIndex].Side = domain.NoneSide
	board.SetMove(g.Moves[limitedIndex])
}

