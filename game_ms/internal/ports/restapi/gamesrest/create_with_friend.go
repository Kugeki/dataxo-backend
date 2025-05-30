package gamesrest

import (
	"dataxo-backend-game-ms/internal/domain"
	"encoding/json"
	"log/slog"
	"net/http"
)

type ModeParams struct {
	MySide int `json:"my_side"`
}

type CreateWithFriendReq struct {
	ModeParams ModeParams `json:"mode_params"`
	ClientID   string     `json:"client_id"`
}

func (r *CreateWithFriendReq) ToDomain() domain.ModeParams {
	return domain.ModeParams{MySide: domain.SideRequest(r.ModeParams.MySide)}
}

type CreateWithFriendResp struct {
	GameID string `json:"game_id"`
}

func (r *CreateWithFriendResp) FromDomain(game *domain.Game) {
	r.GameID = game.ID.String()
}

func (h *Handler) CreateWithFriend() http.HandlerFunc {
	log := h.log
	return func(w http.ResponseWriter, r *http.Request) {
		req := &CreateWithFriendReq{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.responder.RespondError(w, http.StatusBadRequest, err)
			return
		}

		player := domain.PlayerID{ClientID: req.ClientID}
		modeParams := req.ToDomain()
		h.log.Debug("create game", slog.Any("playerID", player))

		g, err := h.gameUC.CreateGame(r.Context(), player, domain.ModeWithFriend, modeParams)
		if err != nil {
			log.Error("uc create game", slog.Any("error", err))
			h.responder.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		resp := &CreateWithFriendResp{}
		resp.FromDomain(g)
		h.responder.Respond(w, http.StatusCreated, resp)
	}
}
