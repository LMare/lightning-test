package handler

import (
	"net/http"
	"fmt"
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
		HtmxResponse(response, "backend/templates/lightning/nodeInfo.html", data)
	} else {
		JsonResponse(response, data)
	}
}

func HandleUpdateNodeInfo(response http.ResponseWriter, request *http.Request) {
	// Parse les données du corps
    err := request.ParseForm()
    if err != nil {
        http.Error(response, "Erreur de parsing", http.StatusBadRequest)
        return
    }

    // Récupère le paramètre "nom"
    alias := request.FormValue("alias")
    fmt.Println("alias reçu : ", alias)
	color := request.FormValue("color")
    fmt.Println("color reçu : ", color)


	// connection info of lnd1
	basePath := config.Load().ProjectPath + "/nodes-storage/lightning-test_lnd1_1"
	authData := service.NewLndClientAuthData(basePath + "/cert/tls.cert", basePath + "/macaroons/admin.macaroon", "localhost:10009");

	err = service.UpdateAliasAndColor(authData, alias, color)
	if err != nil {
		if IsHTMX(request) {
			HtmxMessageKo(response, "Modifications fail.")
		} else {
			// TODO format d'erreur à définir pour unifomisation
		}
	} else {
		if IsHTMX(request) {
			HtmxMessageOk(response, "Modifications successfully applied.")
		} else {
			OkNoContent(response)
		}
	}
}
