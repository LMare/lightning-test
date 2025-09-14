package handler

import (
	"net/http"
	config "github.com/Lmare/lightning-test"
	service "github.com/Lmare/lightning-test/backend/service/lightningService"
)

func HandleNodeInfo(response http.ResponseWriter, request *http.Request) {

	// connection info of lnd1
	basePath := config.Load().ProjectPath + "/nodes-storage/lightning-test_lnd1_1"
	authData := service.NewLndClientAuthData(basePath + "/cert/tls.cert", basePath + "/macaroons/admin.macaroon", "localhost:10009");

	// get the info of the node
	data, err := service.GetUsefullInfo(authData)
	if(err != nil) {
		LogException(err)
	}

	if IsHTMX(request) {
		SetHtmlResponse(response, "backend/templates/lightning/nodeInfo.html", data)
	} else {
		SetJsonResponse(response, data)
	}
}

func HandleUpdateNodeInfo(response http.ResponseWriter, request *http.Request) {

}
