package gamesrest

import (
	"context"
	"dataxo-backend-game-ms/internal/domain"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	"log/slog"
	"net/http"
)

type WsMuxReq struct {
	Type      string          `json:"type"`
	RequestID string          `json:"request_id"`
	ClientID  string          `json:"client_id"`
	Message   json.RawMessage `json:"message"`
}

type WsSideResponse struct {
	Type          string `json:"type"`
	ResponseForID string `json:"response_for_id"`

	Side domain.Side `json:"side"`
}

type WsCreateResponse struct {
	Type   string `json:"type"`
	GameID string `json:"game_id"`
	Side   int    `json:"side"`
}

var CreateMessageType = "create_message"
var SideMessageType = "side_message"

func (h *Handler) WsMux() http.HandlerFunc {
	h.wsHandler.HandleConnect(func(session *melody.Session) {
		h.log.Debug("WebSocket Connected",
			slog.String("remote_addr", session.RemoteAddr().String()),
		)
		ctx := session.Request.Context()

		gameIDParam := chi.URLParam(session.Request, "game_id")

		gameID, err := uuid.Parse(gameIDParam)
		if err != nil && gameIDParam != "create" {
			h.log.Warn("can't parse game id param to uuid",
				slog.String("handler", "WsMux"),
				slog.String("param", gameIDParam),
				slog.Any("error", err),
			)
			h.RespondErrorWsAndClose(session, "", err, h.log)
			return
		}

		clientID := chi.URLParam(session.Request, "client_id")
		session.Set("client_id", clientID)

		// TODO: redo this, it's just workaround
		if gameIDParam == "create" {
			player := domain.PlayerID{
				RemoteAddr: h.GetRemoteAddr(session.Request),
				ClientID:   clientID,
			}
			h.log.Debug("create game", slog.Any("playerID", player))

			g, err := h.gameUC.CreateGame(ctx, player, domain.ModeWithFriend,
				domain.ModeParams{MySide: domain.RandomSideRequest})
			if err != nil {
				h.log.Error("uc create game", slog.Any("error", err))
				h.RespondErrorWsAndClose(session, "", err, h.log)
				return
			}

			side, err := h.WsGetSide(ctx, session, g.ID)
			if err != nil {
				h.log.Error("uc create game get side", slog.Any("error", err))
				h.RespondErrorWsAndClose(session, "", err, h.log)
				return
			}

			resp := WsCreateResponse{
				Type:   CreateMessageType,
				GameID: g.ID.String(),
				Side:   int(side),
			}
			h.wsResponder.RespondWs(session, resp)

			gameID = g.ID
		}

		_, err = h.gameUC.GetGame(ctx, gameID)
		if err != nil {
			h.RespondErrorWsAndClose(session, "", err, h.log)
		}

		session.Set("game_id", gameID)
	})

	h.wsHandler.HandleMessage(func(session *melody.Session, bytes []byte) {
		req := &WsMuxReq{}
		if err := json.Unmarshal(bytes, req); err != nil {
			h.wsResponder.RespondErrorWs(session, err)
			return
		}

		h.log.Debug("WebSocket Mux Request",
			slog.String("json", string(bytes)),
			slog.Any("struct", req),
		)

		gameID, err := h.WsGameIDFromSession(session)
		if err != nil {
			h.wsResponder.RespondErrorWs(session, err)
			return
		}

		if req.ClientID != "" {
			session.Set("client_id", req.ClientID)
		}

		switch req.Type {
		case "presence":
			h.WsPresence(session, req.RequestID, gameID, req.Message)
		case "game":
			h.WsGame(session, req.RequestID, gameID, req.Message)
		case "state":
			h.WsState(session, req.RequestID, gameID, req.Message)
		case "side":
			h.WsSide(session, req.RequestID, gameID, req.Message)
		default:
			h.WsRespondErrorWithID(session, ErrWrongMessageType, req.RequestID)
		}
	})

	h.wsHandler.HandleDisconnect(func(session *melody.Session) {
		h.log.Debug("WebSocket Disconnected",
			slog.String("remote_addr", session.RemoteAddr().String()),
		)
	})

	return func(w http.ResponseWriter, r *http.Request) {
		err := h.wsHandler.HandleRequest(w, r)
		if err != nil {
			return
		}
	}
}

func (h *Handler) RespondErrorWsAndClose(session *melody.Session, requestID string, err error, log *slog.Logger) {
	h.WsRespondErrorWithID(session, err, requestID)
	closeErr := session.Close()
	if closeErr != nil {
		log.Error("can't close WebSocket session", slog.Any("error", closeErr))
	}
}

func (h *Handler) WsGameIDFromSession(session *melody.Session) (uuid.UUID, error) {
	gameIDValue, ok := session.Get("game_id")
	if !ok {
		h.log.Error("can't get game id from session")
		return uuid.UUID{}, ErrCantGetGameIDFromSession
	}

	gameID, ok := gameIDValue.(uuid.UUID)
	if !ok {
		h.log.Error("can't convert game id value to uuid")
		return uuid.UUID{}, ErrCantConvertGameIDValueToUUID
	}

	return gameID, nil
}

func (h *Handler) WsGetSide(ctx context.Context, session *melody.Session, gameID uuid.UUID) (domain.Side, error) {
	sideValue, ok := session.Get("side")
	if !ok {
		return h.wsGetSideFromUC(ctx, session, gameID)
	}

	side, ok := sideValue.(domain.Side)
	if !ok || side == domain.NoneSide {
		return h.wsGetSideFromUC(ctx, session, gameID)
	}

	return side, nil
}

func (h *Handler) wsGetSideFromUC(ctx context.Context, session *melody.Session, gameID uuid.UUID) (domain.Side, error) {
	playerID := h.WsGetPlayerID(session)
	h.log.Debug("side from uc", slog.Any("playerID", playerID))

	side, err := h.gameUC.GetSide(ctx, gameID, playerID)
	if err != nil {
		return domain.NoneSide, err
	}

	session.Set("side", side)
	return side, nil
}

func (h *Handler) WsSendSide(ctx context.Context, session *melody.Session, requestID string, gameID uuid.UUID) {
	side, err := h.WsGetSide(ctx, session, gameID)
	if err != nil {
		side = domain.NoneSide
	}

	h.wsResponder.RespondWs(session, WsSideResponse{
		Type:          SideMessageType,
		Side:          side,
		ResponseForID: requestID,
	})
}

func (h *Handler) WsGetPlayerID(session *melody.Session) domain.PlayerID {
	remoteAddr := h.GetRemoteAddr(session.Request) // todo: session.RemoteAddr().String()?

	clientIDValue, ok := session.Get("client_id")
	if !ok {
		clientIDValue = session.Request.Header.Get("Client-Id")
	}
	clientID, _ := clientIDValue.(string)

	return domain.PlayerID{
		RemoteAddr: remoteAddr,
		ClientID:   clientID,
	}
}

type WsError struct {
	Error         string `json:"error"`
	ResponseForID string `json:"response_for_id,omitempty"`
	NeedReSync    bool   `json:"need_re_sync"`
}

func (h *Handler) WsRespondErrorWithID(session *melody.Session, err error, requestID string) {
	h.wsResponder.RespondWs(session, &WsError{
		Error:         err.Error(),
		ResponseForID: requestID,
		NeedReSync:    domain.IsNeedReSync(err),
	})
}
