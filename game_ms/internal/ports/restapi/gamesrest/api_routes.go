package gamesrest

import "github.com/go-chi/chi/v5"

func (h *Handler) SetupRoutes(r chi.Router) {
	r.Post("/api/v1/games/modes/with-friend", h.CreateWithFriend())
	r.Get("/api/v1/games/{game_id}/state", h.GetState())
	r.Handle("/api/v1/games/{game_id}", h.WsGame())
}
