package handler

import (
	"fmt"
	"net/http"
	"strings"

	Config "github.com/Lmare/lightning-playground"
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus metrics registry (default)
var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"handler", "method"},
	)
)

// Register Prometheus metrics in init
func init() {
	prometheus.MustRegister(requestCounter)
}

type route struct {
	path     string
	handlers map[string]func(http.ResponseWriter, *http.Request)
}

type Router struct {
	routes map[string]*route
}

func GetRouter() *Router {
	router := Router{routes: make(map[string]*route)}
	router.add("/", http.MethodGet, HandleRoot)

	router.add("/lightning/alias", http.MethodPut, handleUpdateNodeAlias)
	router.add("/lightning/channel", http.MethodPost, handleOpenChannel)
	router.add("/lightning/invoice", http.MethodPost, handleCreateInvoice)
	router.add("/lightning/nodes", http.MethodGet, handleListOfNodes)
	router.add("/lightning/nodeInfo", http.MethodGet, handleNodeInfo)
	router.add("/lightning/paiment", http.MethodPost, handleMakePaiment)
	router.add("/lightning/peer", http.MethodPost, handleAddPeer)
	router.add("/lightning/uri", http.MethodGet, handleShowUri)

	router.add("/users", http.MethodGet, HandleListPersonne)

	router.add("/stream-event", http.MethodGet, handleStreamEvent)

	router.add("/health", http.MethodGet, handleHealth)
	router.add("/ready", http.MethodGet, handleReady)
	router.add("/version", http.MethodGet, handleVersion)
	router.add("/metrics", http.MethodGet, handleMetrics)

	return &router
}

// wrapHandlerWithMetrics adds Prometheus instrumentation around the handler
func wrapHandlerWithMetrics(handlerName string, h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCounter.WithLabelValues(handlerName, r.Method).Inc()
		h(w, r)
	}
}

func (router *Router) add(path string, verbe string, callback func(http.ResponseWriter, *http.Request)) {
	if _, exist := router.routes[path]; !exist {
		router.routes[path] = &route{path: path, handlers: make(map[string]func(http.ResponseWriter, *http.Request))}
	}
	// Instruments all handlers except /metrics (avoids double counting on Prometheus scrape)
	handlerName := path
	if path == "/metrics" {
		router.routes[path].handlers[verbe] = callback
	} else {
		router.routes[path].handlers[verbe] = wrapHandlerWithMetrics(handlerName, callback)
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r, ok := router.routes[req.URL.Path]; ok {
		if handler, ok := r.handlers[req.Method]; ok {
			handler(w, req)
			return
		}
		requestCounter.WithLabelValues(req.URL.Path, req.Method).Inc()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, req)
}

/**
 * Handler for / : provides a simple interface to call other routes
 */
func HandleRoot(response http.ResponseWriter, request *http.Request) {

	html := `
	<h1>Lightning-playground</h1>
	<p>
		Bienvenue sur le backend de cette application!
	</p>
	<p>
		Voici la liste des routes déservies :
	</p>
	<ul>`
	for _, r := range GetRouter().routes {
		html = html + "<li><a href=\"" + r.path + "\">" + r.path + "</a> ("
		// Extract methods
		methods := []string{}
		for method := range r.handlers {
			methods = append(methods, method)
		}
		html += strings.Join(methods, ", ") + ")</li>"
	}
	html = html + "</ul>"

	fmt.Fprintf(response, "%s", html)
}

// handleVersion returns the current version of the application
func handleVersion(w http.ResponseWriter, r *http.Request) {
	if IsHTMX(r) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<p>" + Config.Load().Version + "</p>"))
	} else {
		jsonResponse(w, Config.Load().Version)
	}
}
