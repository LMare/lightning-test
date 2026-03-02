package main

import (
	"net/http"

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

func initRouter(handlers *handlers) *Router {
	router := Router{routes: make(map[string]*route)}

	router.add("/lightning/alias", http.MethodPut, handlers.lightning.HandleUpdateNodeAlias)
	router.add("/lightning/channel", http.MethodPost, handlers.lightning.HandleOpenChannel)
	router.add("/lightning/invoice", http.MethodPost, handlers.lightning.HandleCreateInvoice)
	router.add("/lightning/nodes", http.MethodGet, handlers.lightning.HandleListOfNodes)
	router.add("/lightning/nodeInfo", http.MethodGet, handlers.lightning.HandleNodeInfo)
	router.add("/lightning/paiment", http.MethodPost, handlers.lightning.HandleMakePayment)
	router.add("/lightning/peer", http.MethodPost, handlers.lightning.HandleAddPeer)
	router.add("/lightning/uri", http.MethodGet, handlers.lightning.HandleShowUri)

	router.add("/users", http.MethodGet, handlers.user.HandleListUsers)

	router.add("/stream-event", http.MethodGet, handlers.streamEvent.HandleStreamEvent)

	router.add("/health", http.MethodGet, handlers.probe.HandleHealth)
	router.add("/ready", http.MethodGet, handlers.probe.HandleReady)
	router.add("/version", http.MethodGet, handlers.version.HandleVersion)
	router.add("/metrics", http.MethodGet, handlers.probe.HandleMetrics)

	return &router
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
		router.routes[path].handlers[verbe] = router.wrapHandlerWithMetrics(handlerName, callback)
	}
}

// wrapHandlerWithMetrics adds Prometheus instrumentation around the handler
func (router *Router) wrapHandlerWithMetrics(handlerName string, h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCounter.WithLabelValues(handlerName, r.Method).Inc()
		h(w, r)
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
