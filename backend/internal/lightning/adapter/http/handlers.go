package lightninghttp

import (
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Lmare/lightning-playground/backend/internal/httpapi"
	lndgrpc "github.com/Lmare/lightning-playground/backend/internal/lightning/adapter/grpc"
	"github.com/Lmare/lightning-playground/backend/internal/lightning/application"
)

func NewHandler(lightningInfoService *lndgrpc.LightningInfoService, channelService *lndgrpc.ChannelService, nodeService *application.NodeService) *Handler {
	return &Handler{
		lightningInfoService: lightningInfoService,
		channelService:       channelService,
		nodeService:          nodeService,
	}
}

type Handler struct {
	lightningInfoService *lndgrpc.LightningInfoService
	channelService       *lndgrpc.ChannelService
	nodeService          *application.NodeService
}

// get the list of node
func (h *Handler) HandleListOfNodes(response http.ResponseWriter, request *http.Request) {

	descriptors, err := h.nodeService.ListOfNodes()
	if err != nil {
		httpapi.Fail(response, request, "Error on listing nodes", err)
		return
	}

	nodes := h.lightningInfoService.GetListOfNode(descriptors)
	if httpapi.IsHTMX(request) {
		httpapi.HTMXResponse(response, "lightning/nodes.html", nodes)
	} else {
		httpapi.JSONResponse(response, nodes)
	}
}

// check if the format of the id is valid
func isVaildFormatOfId(id string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return re.MatchString(id)
}

// get the URI of the node
func (h *Handler) HandleShowUri(response http.ResponseWriter, request *http.Request) {
	// node id parameter
	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		httpapi.FailCheck(response, request, "Pas d'id transmis", nil)
		return
	}
	// LND client connection details for this node
	authData := h.nodeService.GetLndClientAuthData(id)

	// get the uri of the node
	uri, err := h.lightningInfoService.GetFirstUri(authData)
	if err != nil {
		httpapi.Fail(response, request, "Echec de la communication avec le noeud LND", err)
		return
	}

	// Render
	if httpapi.IsHTMX(request) {
		funcMap := template.FuncMap{"truncateUri": truncateUri}
		httpapi.HTMXResponseWithFuncs(response, "lightning/uri.html", uri, funcMap)
	} else {
		httpapi.JSONResponse(response, uri)
	}

}

// reduce an uri
func truncateUri(s string, n int) string {
	at := strings.Index(s, "@")
	if at == -1 || at < 2*n {
		return s // no @ or too short to truncate
	}

	start := s[:n]
	end := s[at-n : at]
	host := s[at:] // includes @

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
func (h *Handler) HandleNodeInfo(response http.ResponseWriter, request *http.Request) {
	// node id parameter
	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		httpapi.FailCheck(response, request, "Pas d'id transmis", nil)
		return
	}
	// LND client connection details for this node
	authData := h.nodeService.GetLndClientAuthData(id)

	// get the info of the node
	data, err := h.lightningInfoService.GetUsefullInfo(authData)
	if err != nil {
		httpapi.Fail(response, request, "Echec de la communication avec le noeud LND", err)
		return
	}

	// Render
	if httpapi.IsHTMX(request) {
		httpapi.HTMXResponse(response, "lightning/nodeInfo.html", data)
	} else {
		httpapi.JSONResponse(response, data)
	}
}

// Create a connexion to a new Peer
func (h *Handler) HandleAddPeer(response http.ResponseWriter, request *http.Request) {
	// Parse request body (form)
	err := request.ParseForm()
	if err != nil {
		http.Error(response, "Erreur de parsing", http.StatusBadRequest)
		return
	}

	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		httpapi.FailCheck(response, request, "Pas d'id transmis", nil)
		return
	}
	uri := request.FormValue("uri")
	// Get Data to connect lnd
	authData := h.nodeService.GetLndClientAuthData(id)

	// Add the pair
	err = h.channelService.AddPeer(authData, uri)
	if err != nil {
		httpapi.Fail(response, request, "Fail to add the peer.", err)
		return
	}

	if httpapi.IsHTMX(request) {
		httpapi.HTMXMessageOk(response, "Peer successfully added.")
	} else {
		httpapi.OkNoContent(response)
	}
}

// create a channel
func (h *Handler) HandleOpenChannel(response http.ResponseWriter, request *http.Request) {

	// Parse request body (form)
	err := request.ParseForm()
	if err != nil {
		httpapi.FailCheck(response, request, "Parse error", err)
		return
	}

	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		httpapi.FailCheck(response, request, "Pas d'id transmis", nil)
		return
	}
	pubKey := request.FormValue("pubKey")
	amountStr := request.FormValue("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		httpapi.FailCheck(response, request, "Amount value incorrect", err)
		return
	}

	// Get Data to connect lnd
	authData := h.nodeService.GetLndClientAuthData(id)

	// create the channel
	err = h.channelService.OpenChannel(authData, pubKey, amount)
	if err != nil {
		httpapi.Fail(response, request, "Fail to create the channel.", err)
		return
	}

	if httpapi.IsHTMX(request) {
		httpapi.HTMXMessageOk(response, "Channel successfully created.")
	} else {
		httpapi.OkNoContent(response)
	}
}

// Create an invoice
func (h *Handler) HandleCreateInvoice(response http.ResponseWriter, request *http.Request) {

	// Parse request body (form)
	err := request.ParseForm()
	if err != nil {
		httpapi.FailCheck(response, request, "Parse error", err)
		return
	}

	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		httpapi.FailCheck(response, request, "Pas d'id transmis", nil)
		return
	}
	memo := request.FormValue("memo")
	amountStr := request.FormValue("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		httpapi.FailCheck(response, request, "Amount value incorrect", err)
		return
	}

	// Get Data to connect lnd
	authData := h.nodeService.GetLndClientAuthData(id)

	// create the invoice
	p, err := h.channelService.CreateQuickInvoice(authData, memo, amount)
	if err != nil {
		httpapi.Fail(response, request, "Fail to create the invoice.", err)
		return
	}

	// Render
	if httpapi.IsHTMX(request) {
		funcMap := template.FuncMap{"truncate": truncate}
		httpapi.HTMXResponseWithFuncs(response, "lightning/paymentRequest.html", p, funcMap)
	} else {
		httpapi.JSONResponse(response, p)
	}

}

// Make the payment
func (h *Handler) HandleMakePayment(response http.ResponseWriter, request *http.Request) {
	// Parse request body (form)
	err := request.ParseForm()
	if err != nil {
		httpapi.FailCheck(response, request, "Error parsing form", err)
		return
	}

	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		httpapi.FailCheck(response, request, "Missing ID", nil)
		return
	}
	paymentRequest := request.FormValue("paymentRequest")
	if paymentRequest == "" {
		httpapi.FailCheck(response, request, "Missing payment request", nil)
		return
	}

	// Get Data to connect lnd
	authData := h.nodeService.GetLndClientAuthData(id)

	// pay the invoice (async updates via SSE)
	err = h.channelService.MakePayment(authData, paymentRequest)
	if err != nil {
		httpapi.Fail(response, request, "Fail to pay the invoice.", err)
		return
	}

	// TODO have a streamId to return ?
	if httpapi.IsHTMX(request) {
		httpapi.HTMXMessageOk(response, "Processing in progress")
	} else {
		httpapi.OkNoContent(response)
	}

}

// Update name of the node & color
func (h *Handler) HandleUpdateNodeAlias(response http.ResponseWriter, request *http.Request) {
	// Parse request body (form)
	err := request.ParseForm()
	if err != nil {
		httpapi.FailCheck(response, request, "Parse error", err)
		return
	}

	alias := request.FormValue("alias")
	if alias == "" {
		httpapi.FailCheck(response, request, "Missing alias", nil)
		return
	}
	color := request.FormValue("color")
	if color == "" {
		httpapi.FailCheck(response, request, "Missing color", nil)
		return
	}

	// node id parameter
	id := request.FormValue("id")
	if !isVaildFormatOfId(id) {
		httpapi.FailCheck(response, request, "Missing ID", nil)
		return
	}

	// LND connection info for this node
	authData := h.nodeService.GetLndClientAuthData(id)

	err = h.lightningInfoService.UpdateAliasAndColor(authData, alias, color)
	if err != nil {
		httpapi.Fail(response, request, "Modifications fail.", err)
	}

	if httpapi.IsHTMX(request) {
		httpapi.HTMXMessageOk(response, "Modifications successfully applied.")
	} else {
		httpapi.OkNoContent(response)
	}

}
