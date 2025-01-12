package gamesrest

import (
	"github.com/google/uuid"
	"github.com/olahol/melody"
)

func (h *Handler) WsSide(session *melody.Session, requestID string, gameID uuid.UUID, bytes []byte) {
	h.WsSendSide(session.Request.Context(), session, requestID, gameID)
}
