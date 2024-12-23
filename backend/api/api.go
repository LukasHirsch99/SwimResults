package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"swimresults-backend/internal/repository"
	"time"
)

type Api struct {
	router *http.ServeMux
	repo   *repository.Queries
	logger *slog.Logger
}

func New(repo *repository.Queries, logger *slog.Logger) *Api {
	return &Api{
		router: http.NewServeMux(),
		repo:   repo,
		logger: logger,
	}
}

func (a *Api) Start(ctx context.Context) error {
	if a.repo == nil {
		return fmt.Errorf("repo empty")
	}

  a.loadRoutes()

	server := http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	done := make(chan struct{})
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("failed to listen and serve", slog.Any("error", err))
		}
		close(done)
	}()

	a.logger.Info("Server listening", slog.String("addr", ":8080"))
	select {
	case <-done:
		break
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		server.Shutdown(ctx)
		cancel()
	}

	return nil
}
