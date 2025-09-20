package handler

import (
	"net/http"
	"fmt"
	config "github.com/Lmare/lightning-test"
	lightningService "github.com/Lmare/lightning-test/backend/service/lightningService"
	nodeService "github.com/Lmare/lightning-test/backend/service/nodeService"
)

func HandleNodeInfo(response http.ResponseWriter, request *http.Request) {

	authData, err := nodeService.GetLndClientAuthData(1)
	if(err != nil) {
		fail(response, request, "Info transmisent incorrectes", err)
		return
	}

	// get the info of the node
	data, err := lightningService.GetUsefullInfo(authData)
	if(err != nil) {
		fail(response, request, "Echec de la communication avec le noeud LND", err)
		return
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
	authData := lightningService.NewLndClientAuthData(basePath + "/cert/tls.cert", basePath + "/macaroons/admin.macaroon", "localhost:10009");

	err = lightningService.UpdateAliasAndColor(authData, alias, color)
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
