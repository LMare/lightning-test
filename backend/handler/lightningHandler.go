package handler

import (
	"net/http"
	"fmt"

	service "github.com/Lmare/lightning-test/backend/service/lightningService"
	//"github.com/Lmare/lightning-test/backend/templates/personView"
)

func HandleNodeInfo(response http.ResponseWriter, request *http.Request) {

	basePath := "/home/louis/Documents/Dev/lightning-test/nodes-storage/lightning-test_lnd1_1"

	data, err := service.GetUsefullInfo(service.NewLndClientAuthData(basePath + "/cert/tls.cert", basePath + "/macaroons/admin.macaroon", "localhost:10009"))
	if(err != nil) {
		fmt.Println("Une erreur est survenue : ", err)
	}


	if IsHTMX(request) {
		SetHtmlResponse(response, "backend/templates/lightning/nodeInfo.html", data)
	} else {
		SetJsonResponse(response, data)
	}
}
