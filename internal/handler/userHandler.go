package handler

import (
	"net/http"

	"github.com/Lmare/lightning-test/internal/service"
)

func HandleListPersonne(response http.ResponseWriter, request *http.Request) {

	users := service.ListUsers()

	SetJsonResponse(response, users)
}
