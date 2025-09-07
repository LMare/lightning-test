package handler

import (
	"net/http"

	"github.com/Lmare/lightning-test/internal/service"
	"github.com/Lmare/lightning-test/internal/templates/personView"
)

func HandleListPersonne(response http.ResponseWriter, request *http.Request) {

	users := service.ListUsers()
	if IsHTMX(request) {
		vo := personView.ViewObject(users)
		SetHtmlResponse(response, "internal/templates/personView/user.html", vo)
	} else {
		SetJsonResponse(response, users)
	}
}
