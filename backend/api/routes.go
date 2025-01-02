package api

import (
	"net/http"
	"swimresults-backend/api/handler"
)

func (a *Api) loadRoutes() {
	swimresults := handler.New(a.logger, a.repo)

  a.router.Handle("GET /", http.HandlerFunc(swimresults.GetHome))
  a.router.Handle("GET /meets/", http.HandlerFunc(swimresults.GetMeets))
  a.router.Handle("GET /swimmers/", http.HandlerFunc(swimresults.GetSwimmers))
  a.router.Handle("GET /clubs/", http.HandlerFunc(swimresults.GetClubs))
}
