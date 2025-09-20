package handler

import (
	"net/http"
	"fmt"
	"strconv"

	lightningService "github.com/Lmare/lightning-test/backend/service/lightningService"
	nodeService "github.com/Lmare/lightning-test/backend/service/nodeService"
)


func HandleListOfNodes(response http.ResponseWriter, request *http.Request) {

	descriptors, err := nodeService.ListOfNodes()
	if(err != nil) {
		fail(response, request, "Info transmisent incorrectes", err)
		return
	}

	nodes := lightningService.GetListOfNode(descriptors)
	if IsHTMX(request) {
		HtmxResponse(response, "lightning/nodes.html", nodes)
	} else {
		JsonResponse(response, nodes)
	}
}


// get the info of one Node
func HandleNodeInfo(response http.ResponseWriter, request *http.Request) {
	// paramètre de la node
	idStr := request.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fail(response, request, "Pas d'id transmis", err)
		return
	}
	// récupération info de connexion à la node
	authData, err := nodeService.GetLndClientAuthData(id)
	if(err != nil) {
		fail(response, request, "Node inexistante", err)
		return
	}

	// get the info of the node
	data, err := lightningService.GetUsefullInfo(authData)
	if(err != nil) {
		fail(response, request, "Echec de la communication avec le noeud LND", err)
		return
	}

	// Render
	if IsHTMX(request) {
		HtmxResponse(response, "lightning/nodeInfo.html", data)
	} else {
		JsonResponse(response, data)
	}
}


// Update name of the node & color
// TODO : update lnd to have gRPC methode to do that
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
	authData, err := nodeService.GetLndClientAuthData(1)
	if(err != nil) {
		fail(response, request, "Info transmisent incorrectes", err)
		return
	}


	err = lightningService.UpdateAliasAndColor(authData, alias, color)
	if err != nil {
		fail(response, request, "Modifications fail.", err)
	}

	if IsHTMX(request) {
		HtmxMessageOk(response, "Modifications successfully applied.")
	} else {
		OkNoContent(response)
	}

}
