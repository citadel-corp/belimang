package merchants

import (
	"fmt"
	"net/http"

	"github.com/citadel-corp/belimang/internal/common/request"
	"github.com/citadel-corp/belimang/internal/common/response"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateMerchantPayload

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}

	userResp, err := h.service.Create(r.Context(), req)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusCreated, response.ResponseBody{
		Message: "Merchant created successfully",
		Data:    userResp,
	})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var req ListMerchantsPayload

	newSchema := schema.NewDecoder()
	newSchema.IgnoreUnknownKeys(true)

	if err := newSchema.Decode(&req, r.URL.Query()); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{})
		return
	}

	err := req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}

	merchantsResp, err := h.service.List(r.Context(), req)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "Merchants fetched successfully",
		Data:    merchantsResp,
	})
}

func (h *Handler) ListByDistance(w http.ResponseWriter, r *http.Request) {
	var req ListMerchantsByDistancePayload

	newSchema := schema.NewDecoder()
	newSchema.IgnoreUnknownKeys(true)

	if err := newSchema.Decode(&req, r.URL.Query()); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{})
		return
	}

	req.Lat = mux.Vars(r)["lat"]
	req.Lng = mux.Vars(r)["long"]

	fmt.Println("lat, lng: ", req.Lat, req.Lng)

	err := req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}

	merchantsResp, err := h.service.ListByDistance(r.Context(), req)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "Merchants fetched successfully",
		Data:    merchantsResp,
	})
}
