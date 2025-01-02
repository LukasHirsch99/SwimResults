package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"swimresults-backend/internal/repository"
)

type SwimResults struct {
	logger *slog.Logger
	repo   *repository.Queries
}

func New(logger *slog.Logger, repo *repository.Queries) *SwimResults {
	return &SwimResults{
		repo:   repo,
		logger: logger,
	}
}

type indexPage struct {
	SwimmerIds []int32
	Total      int64
}

type errorPage struct {
	ErrorMessage string
}

func (h *SwimResults) GetHome(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode("Swim Results Api")
}

func (h *SwimResults) GetMeets(w http.ResponseWriter, r *http.Request) {
	meets, err := h.repo.GetMeets(r.Context())
	if err != nil {
		h.logger.Error("failed to get meets", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meets)
}

func (h *SwimResults) GetClubs(w http.ResponseWriter, r *http.Request) {
	clubs, err := h.repo.GetClubs(r.Context())
	if err != nil {
		h.logger.Error("failed to get clubs", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clubs)
}

func (h *SwimResults) GetSwimmers(w http.ResponseWriter, r *http.Request) {
	swimmerIds, err := h.repo.GetSwimmerIds(r.Context())
	if err != nil {
		h.logger.Error("failed to find swimmerids", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := indexPage{
		SwimmerIds: swimmerIds,
		Total:      int64(len(swimmerIds)),
	}
	json.NewEncoder(w).Encode(response)
}
