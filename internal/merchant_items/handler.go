package merchantitems

import (
	"net/http"

	"github.com/citadel-corp/belimang/internal/common/request"
	"github.com/citadel-corp/belimang/internal/common/response"
	"github.com/citadel-corp/belimang/internal/merchants"
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
	var req CreateMerchantItemPayload

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	req.MerchantID = mux.Vars(r)["merchantId"]

	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}

	itemResp, err := h.service.Create(r.Context(), req)
	if err == merchants.ErrMerchantNotFound {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Merchant not found",
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
	response.JSON(w, http.StatusCreated, response.ResponseBody{
		Message: "Merchant item created successfully",
		Data:    itemResp,
	})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var req ListMerchantItemsPayload

	newSchema := schema.NewDecoder()
	newSchema.IgnoreUnknownKeys(true)

	if err := newSchema.Decode(&req, r.URL.Query()); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{})
		return
	}

	req.MerchantUID = mux.Vars(r)["merchantId"]

	err := req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}

	itemResp, err := h.service.List(r.Context(), req)
	if err == merchants.ErrMerchantNotFound {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Merchant not found",
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
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "Merchant items fetched successfully",
		Data:    itemResp,
	})
}
