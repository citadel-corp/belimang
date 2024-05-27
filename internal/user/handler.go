package user

import (
	"errors"
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

var (
	requestCreate CreateUserPayload
)

func (h *Handler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	err := request.DecodeJSON(w, r, &requestCreate)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	requestCreate.UserType = Admin

	err = requestCreate.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}

	h.CreateUser(w, r)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	userResp, err := h.service.Create(r.Context(), requestCreate)
	if errors.Is(err, ErrUserAlreadyExists) {
		response.JSON(w, http.StatusConflict, response.ResponseBody{
			Message: "User already exists",
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
		Message: "User registered successfully",
		Data:    userResp,
	})
}
