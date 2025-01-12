package gamesrest

import (
	"dataxo-backend-game-ms/internal/ports/restapi"
	"dataxo-backend-game-ms/pkg/slogdiscard"
	"github.com/olahol/melody"
	"log/slog"
	"net"
	"net/http"
)

type Handler struct {
	log         *slog.Logger
	gameUC      GameUsecase
	responder   restapi.Responder
	wsResponder restapi.WsResponder

	wsHandler *melody.Melody
}

func New(log *slog.Logger, gameUC GameUsecase, responder restapi.Responder, wsResponder restapi.WsResponder) *Handler {
	log = slogdiscard.LoggerIfNil(log)
	return &Handler{log: log, gameUC: gameUC, responder: responder, wsResponder: wsResponder, wsHandler: melody.New()}
}

func (h *Handler) GetRemoteAddr(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		h.log.Error("get remote addr: split host port", slog.Any("error", err))
	}
	return host
}
