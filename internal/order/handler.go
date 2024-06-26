package order

import (
	"errors"
	"net/http"

	"github.com/citadel-corp/belimang/internal/common/haversine"
	"github.com/citadel-corp/belimang/internal/common/jwt"
	"github.com/citadel-corp/belimang/internal/common/middleware"
	"github.com/citadel-corp/belimang/internal/common/request"
	"github.com/citadel-corp/belimang/internal/common/response"
	"github.com/gorilla/schema"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CalculateEstimate(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
		return
	}
	var req CalculateOrderEstimateRequest

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	res, err := h.service.CalculateEstimate(r.Context(), req, userID)
	if errors.Is(err, ErrValidationFailed) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, haversine.ErrDistanceTooFar) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrSomeMerchantNotFound) {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Not found",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrSomeItemNotFound) {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Not found",
			Error:   err.Error(),
		})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
		return
	}
	var req CreateOrderRequest

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	res, err := h.service.CreateOrder(r.Context(), req, userID)
	if errors.Is(err, ErrValidationFailed) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrCalculatedEstimateNotFound) {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Not found",
			Error:   err.Error(),
		})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusCreated, res)
}

func (h *Handler) SearchOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
		return
	}
	var req SearchOrderPayload

	newSchema := schema.NewDecoder()
	newSchema.IgnoreUnknownKeys(true)

	if err := newSchema.Decode(&req, r.URL.Query()); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{})
		return
	}

	orders, err := h.service.SearchOrders(r.Context(), req, userID)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusOK, orders)
}

func getUserID(r *http.Request) (string, error) {
	if authValue, ok := r.Context().Value(middleware.ContextAuthKey{}).(*jwt.UserClaims); ok {
		return authValue.UserUID, nil
	} else {
		return "", errors.New("cannot parse auth value from context")
	}
}
