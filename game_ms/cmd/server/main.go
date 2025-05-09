package main

import (
	"context"
	"dataxo-backend-game-ms/internal/adapters/mapstore"
	"dataxo-backend-game-ms/internal/domain"
	"dataxo-backend-game-ms/internal/ports/restapi"
	"dataxo-backend-game-ms/internal/ports/restapi/gamesrest"
	"dataxo-backend-game-ms/internal/usecases/gameuc"
	"dataxo-backend-game-ms/internal/usecases/gameuc/modes"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lmittmann/tint"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	logLevel = slog.LevelDebug

	restAddr         = ":8080"
	restReadTimeout  = 10 * time.Second
	restWriteTimeout = 10 * time.Second

	shutdownTimeout = 15 * time.Second
)

func main() {
	log := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.StampMilli,
		AddSource:  true,
		Level:      logLevel,
	}))
	slog.SetDefault(log)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	gameRepo := mapstore.NewGameRepo()

	disappearingModeCfg := domain.DisappearingModeConfig{
		PlayerFiguresLimit: 6,
		WinLineLength:      4,
		BoardWidth:         4,
		BoardHeight:        4,
	}

	disappearingMode, err := modes.NewDisappearingMode(disappearingModeCfg, log)
	if err != nil {
		log.Error("can't create disappearing game mode", slog.Any("error", err))
		return
	}

	gameUC := gameuc.New(gameRepo, disappearingMode, log)

	jsonResponder := restapi.NewJsonResponder(log)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RequestLogger(restapi.NewLogFormatter(log, slog.LevelInfo)))
	router.Use(middleware.Recoverer)

	gamesRestHandler := gamesrest.New(log, gameUC, jsonResponder, jsonResponder)
	gamesRestHandler.SetupRoutes(router)

	restOpts := []restapi.Opt{
		restapi.WithAddr(restAddr),
		restapi.WithErrorLog(slog.NewLogLogger(log.Handler(), slog.LevelError)),
		restapi.WithReadTimeout(restReadTimeout),
		restapi.WithWriteTimeout(restWriteTimeout),
	}
	restAPI, err := restapi.New(router, log, restOpts...)
	if err != nil {
		log.Error("can't create rest api", slog.Any("error", err))
		return
	}

	go func() {
		err := restAPI.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("rest api run error", slog.Any("error", err))
			cancel()
		}
	}()

	<-ctx.Done()

	log.Info("graceful shutdown is beginning...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	done := make(chan struct{}, 1)

	go func() {
		err := restAPI.Shutdown(shutdownCtx)
		if err != nil {
			log.Error("rest server shutdown error", slog.Any("error", err))
		}
		close(done)
	}()

	select {
	case <-shutdownCtx.Done():
		log.Error("graceful shutdown has failed", slog.Any("error", shutdownCtx.Err()))
	case <-done:
		log.Info("graceful shutdown has completed successfully")
	}
}
