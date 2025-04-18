package handler

import (
	"net/http"
	"pie_fire_dine/contract"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	resp := contract.NewResponse(http.StatusOK, nil)
	WriteResponse(w, &resp)
}
