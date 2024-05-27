package merchants

import (
	"net/http"

	"github.com/citadel-corp/belimang/internal/common/request"
	"github.com/citadel-corp/belimang/internal/common/response"
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
