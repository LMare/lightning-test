package handler

import (
	"net/http"

	service "github.com/Lmare/lightning-test/backend/service/lightningService"
	//"github.com/Lmare/lightning-test/backend/templates/personView"
)

func HandleLigthningTest(response http.ResponseWriter, request *http.Request) {

	service.Test()
	/*if IsHTMX(request) {
		vo := personView.ViewObject(users)
		SetHtmlResponse(response, "backend/templates/personView/user.html", vo)
	} else {
		SetJsonResponse(response, users)
	}*/
}
