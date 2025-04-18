package server

import (
	"net/http"

	"pie_fire_dine/handler"

	"github.com/gorilla/mux"
)

func Router(pfdHandler handler.MeatSummaryHandler) http.Handler {
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	router.HandleFunc("/health", handler.HealthCheckHandler).Methods(http.MethodGet)
	router.HandleFunc("/{category}/summary", pfdHandler.GetMeatSummary).Methods(http.MethodGet)

	return router
}
