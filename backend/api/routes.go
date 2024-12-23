package api

import (
	"net/http"
	"swimresults-backend/api/handler"
)

func (a *Api) loadRoutes() {
	swimresults := handler.New(a.logger, a.repo)

	a.router.Handle("GET /", http.HandlerFunc(swimresults.Home))
}
