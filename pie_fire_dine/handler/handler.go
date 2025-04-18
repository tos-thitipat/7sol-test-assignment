package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"pie_fire_dine/contract"
	"pie_fire_dine/errs"
	"pie_fire_dine/service"

	"github.com/gorilla/mux"
)

type MeatSummaryHandler interface {
	GetMeatSummary(w http.ResponseWriter, r *http.Request)
}

type meatSummaryHandler struct {
	meatSummaryService service.MeatSummary
}

func NewMeatSummaryHandler(meatSummaryService service.MeatSummary) MeatSummaryHandler {
	return &meatSummaryHandler{
		meatSummaryService: meatSummaryService,
	}
}

func (h *meatSummaryHandler) GetMeatSummary(w http.ResponseWriter, r *http.Request) {
	category, contain := mux.Vars(r)["category"]
	if !contain {
		resp := contract.NewErrorResponse(errs.NewBadRequest())
		WriteResponse(w, &resp)
	}
	response, err := h.meatSummaryService.GetMeatSummary(r.Context(), category)
	if err != nil {
		resp := contract.NewErrorResponse(err)
		WriteResponse(w, &resp)
		return
	}
	resp := contract.NewResponse(http.StatusOK, response)
	WriteResponse(w, &resp)
}

func WriteResponse(w http.ResponseWriter, resp *contract.Response) {
	w.Header().Set("Content-Type", "application/json")

	responseJSON, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		slog.Error("unable to marshal request response")
		return
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write(responseJSON)
}
