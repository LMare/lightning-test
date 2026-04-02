package userhttp

import (
	"net/http"

	"github.com/Lmare/lightning-playground/backend/internal/httpapi"
	"github.com/Lmare/lightning-playground/backend/internal/user/application"
	"github.com/Lmare/lightning-playground/backend/templates/personView"
)

type Handler struct {
	service *application.UserService
}

func NewHandler(service *application.UserService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleListUsers(response http.ResponseWriter, request *http.Request) {

	users, err := h.service.List(request.Context())
	if err != nil {
		httpapi.Fail(response, request, "Error in the user service", err)
		return
	}

	if httpapi.IsHTMX(request) {
		vo := personView.ViewObject(users)
		httpapi.HTMXResponse(response, "personView/user.html", vo)
	} else {
		httpapi.JSONResponse(response, users)
	}
}
