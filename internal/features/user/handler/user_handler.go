package userhandler

import (
	"encoding/json"
	"net/http"

	"github.com/codepnw/stdlib-ticket-system/internal/features/user"
	userusecase "github.com/codepnw/stdlib-ticket-system/internal/features/user/usecase"
	"github.com/codepnw/stdlib-ticket-system/internal/helper"
	"github.com/codepnw/stdlib-ticket-system/pkg/utils"
)

type userHandler struct {
	uc userusecase.UserUsecase
}

func NewUserHandler(uc userusecase.UserUsecase) *userHandler {
	return &userHandler{uc: uc}
}

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req UserCredentials

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := utils.Validate(&req); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.uc.Register(r.Context(), user.User{
		Username:     req.Username,
		HashPassword: req.Password,
	})
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(w, http.StatusCreated, "register successful", data)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req UserCredentials

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := utils.Validate(&req); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.uc.Login(r.Context(), user.User{
		Username:     req.Username,
		HashPassword: req.Password,
	})
	if err != nil {
		helper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(w, http.StatusOK, "login successful", data)
}
