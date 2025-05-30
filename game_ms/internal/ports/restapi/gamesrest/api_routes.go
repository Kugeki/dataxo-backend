package gamesrest

import "github.com/go-chi/chi/v5"

func (h *Handler) SetupRoutes(r chi.Router) {
	r.Post("/api/v1/games/modes/with-friend", h.CreateWithFriend())
	r.Handle("/api/v1/games/{game_id}/{client_id}", h.WsMux())
}
