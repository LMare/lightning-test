package handler

import (
	"net/http"

	"github.com/Lmare/lightning-playground/backend/service/userService"
	"github.com/Lmare/lightning-playground/backend/templates/personView"
)

type UserHandler struct {
	service *userService.UserService
}

func NewUserHandler(service *userService.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) HandleListUsers(response http.ResponseWriter, request *http.Request) {

	users, err := h.service.List(request.Context())
	if err != nil {
		fail(response, request, "Error in the user service", err)
		return
	}

	if IsHTMX(request) {
		vo := personView.ViewObject(users)
		htmxResponse(response, "personView/user.html", vo)
	} else {
		jsonResponse(response, users)
	}
}
