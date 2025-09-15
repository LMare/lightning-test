package handler

import (
	"net/http"

	"github.com/Lmare/lightning-test/backend/service/personneService"
	"github.com/Lmare/lightning-test/backend/templates/personView"
)

func HandleListPersonne(response http.ResponseWriter, request *http.Request) {

	users := personneService.ListUsers()
	if IsHTMX(request) {
		vo := personView.ViewObject(users)
		HtmxResponse(response, "backend/templates/personView/user.html", vo)
	} else {
		JsonResponse(response, users)
	}
}
