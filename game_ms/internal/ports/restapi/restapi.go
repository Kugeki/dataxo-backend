package restapi

import (
	"context"
	"dataxo-backend-game-ms/pkg/slogdiscard"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type RestAPI struct {
	s   *http.Server
	log *slog.Logger
}

func New(r chi.Router, log *slog.Logger, options ...Opt) (*RestAPI, error) {
	s := &http.Server{Handler: r}

	for _, op := range options {
		err := op(s)
		if err != nil {
			return nil, err
		}
	}

	log = slogdiscard.LoggerIfNil(log)
	log = log.With(slog.String("component", "rest api"))
	return &RestAPI{s: s, log: log}, nil
}

func (a *RestAPI) Run() error {
	a.log.Info("listening", slog.String("addr", a.s.Addr))
	return a.s.ListenAndServe()
}

func (a *RestAPI) Shutdown(ctx context.Context) error {
	a.log.Info("shutdown is started")
	return a.s.Shutdown(ctx)
}
