package handler

import (
	"net/http"
	"fmt"
	"strconv"
	"html/template"
	"strings"
	"regexp"

	lightningService "github.com/Lmare/lightning-playground/backend/service/lightningService"
	nodeService "github.com/Lmare/lightning-playground/backend/service/nodeService"
)

// get the list of node
func handleListOfNodes(response http.ResponseWriter, request *http.Request) {

	descriptors, err := nodeService.ListOfNodes()
	if(err != nil) {
		fail(response, request, "Info transmisent incorrectes", err)
		return
	}

	nodes := lightningService.GetListOfNode(descriptors)
	if IsHTMX(request) {
		htmxResponse(response, "lightning/nodes.html", nodes)
	} else {
		jsonResponse(response, nodes)
	}
}

// check if the format of the id is valid
func isVaildFormatOfId(id string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return re.MatchString(id)
}

// get the URI of the node
func handleShowUri(response http.ResponseWriter, request *http.Request) {
	// paramètre de la node
	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		fail(response, request, "Pas d'id transmis", nil)
		return
	}
	// récupération info de connexion à la node
	authData := nodeService.GetLndClientAuthData(id)

	// get the uri of the node
	uri, err := lightningService.GetFirstUri(authData)
	if(err != nil) {
		fail(response, request, "Echec de la communication avec le noeud LND", err)
		return
	}

	// Render
	if IsHTMX(request) {
		funcMap := template.FuncMap{"truncateUri": truncateUri,}
		htmxResponseWithFuncs(response, "lightning/uri.html", uri, funcMap)
	} else {
		jsonResponse(response, uri)
	}

}

// reduce an uri
func truncateUri(s string, n int) string {
    at := strings.Index(s, "@")
    if at == -1 || at < 2*n {
        return s // pas de @ ou trop court pour tronquer
    }

    start := s[:n]
    end := s[at-n : at]
    host := s[at:] // inclut le @

    return start + "..." + end + host
}


// reduce a string
func truncate(s string, n int) string {
	at := len(s)
	if at <= 2*n {
		return s
	}
	start := s[:n]
	end := s[at-n : at]
	return start + "..." + end
}



// get the info of one Node
func handleNodeInfo(response http.ResponseWriter, request *http.Request) {
	// paramètre de la node
	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		fail(response, request, "Pas d'id transmis", nil)
		return
	}
	// récupération info de connexion à la node
	authData := nodeService.GetLndClientAuthData(id)

	// get the info of the node
	data, err := lightningService.GetUsefullInfo(authData)
	if(err != nil) {
		fail(response, request, "Echec de la communication avec le noeud LND", err)
		return
	}

	// Render
	if IsHTMX(request) {
		htmxResponse(response, "lightning/nodeInfo.html", data)
	} else {
		jsonResponse(response, data)
	}
}

// Create a connexion to a new Peer
func handleAddPeer(response http.ResponseWriter, request *http.Request) {
	// Parse les données du corps
    err := request.ParseForm()
    if err != nil {
        http.Error(response, "Erreur de parsing", http.StatusBadRequest)
        return
    }

	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		fail(response, request, "Pas d'id transmis", err)
		return
	}
	uri := request.FormValue("uri")
	// Get Data to connect lnd
	authData := nodeService.GetLndClientAuthData(id)

	// Add the pair
	err = lightningService.AddPeer(authData, uri)
	if err != nil {
		fail(response, request, "Fail to add the peer.", err)
		return
	}

	if IsHTMX(request) {
		htmxMessageOk(response, "Peer successfully added.")
	} else {
		OkNoContent(response)
	}
}


// create a channel
func handleOpenChannel(response http.ResponseWriter, request *http.Request) {

	// Parse les données du corps
    err := request.ParseForm()
    if err != nil {
        http.Error(response, "Erreur de parsing", http.StatusBadRequest)
        return
    }

	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		fail(response, request, "Pas d'id transmis", err)
		return
	}
	pubKey := request.FormValue("pubKey")
	amountStr := request.FormValue("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		fail(response, request, "Amount value incorrect", err)
		return
	}

	// Get Data to connect lnd
	authData := nodeService.GetLndClientAuthData(id)

	// create the channel
	err = lightningService.OpenChannel(authData, pubKey, amount)
	if err != nil {
		fail(response, request, "Fail to create the channel.", err)
		return
	}

	if IsHTMX(request) {
		htmxMessageOk(response, "Channel successfully created.")
	} else {
		OkNoContent(response)
	}
}

// Create an invoice
func handleCreateInvoice(response http.ResponseWriter, request *http.Request) {

	// Parse les données du corps
	err := request.ParseForm()
	if err != nil {
		http.Error(response, "Erreur de parsing", http.StatusBadRequest)
		return
	}

	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		fail(response, request, "Pas d'id transmis", err)
		return
	}
	memo := request.FormValue("memo")
	amountStr := request.FormValue("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		fail(response, request, "Amount value incorrect", err)
		return
	}

	// Get Data to connect lnd
	authData := nodeService.GetLndClientAuthData(id)

	// create the invoice
	p, err := lightningService.CreateQuickInvoice(authData, memo, amount)
	if err != nil {
		fail(response, request, "Fail to create the invoice.", err)
		return
	}

	// Render
	if IsHTMX(request) {
		funcMap := template.FuncMap{"truncate": truncate,}
		htmxResponseWithFuncs(response, "lightning/paymentRequest.html", p, funcMap)
	} else {
		jsonResponse(response, p)
	}

}

// Make the payment
func handleMakePaiment(response http.ResponseWriter, request *http.Request) {
	// Parse les données du corps
	err := request.ParseForm()
	if err != nil {
		http.Error(response, "Erreur de parsing", http.StatusBadRequest)
		return
	}

	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		fail(response, request, "Pas d'id transmis", err)
		return
	}
	paymentRequest := request.FormValue("paymentRequest")

	// Get Data to connect lnd
	authData := nodeService.GetLndClientAuthData(id)

	// create the invoice
	err = lightningService.MakePaiment(authData, paymentRequest)
	if err != nil {
		fail(response, request, "Fail to pay the invoice.", err)
		return
	}

	// TODO have a streamId to return ?
	if IsHTMX(request) {
		htmxMessageOk(response, "Processing in progress")
	} else {
		OkNoContent(response)
	}

}



// Update name of the node & color
func handleUpdateNodeAlias(response http.ResponseWriter, request *http.Request) {
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

	// paramètre de la node
	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		fail(response, request, "Pas d'id transmis", err)
		return
	}

	// connection info of lnd1
	authData := nodeService.GetLndClientAuthData(id)

	err = lightningService.UpdateAliasAndColor(authData, alias, color)
	if err != nil {
		fail(response, request, "Modifications fail.", err)
	}

	if IsHTMX(request) {
		htmxMessageOk(response, "Modifications successfully applied.")
	} else {
		OkNoContent(response)
	}

}
