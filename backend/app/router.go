package app

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type route struct {
	path     string
	handlers map[string]func(http.ResponseWriter, *http.Request)
}

type router struct {
	routes   map[string]*route
	metrics  *prometheus.CounterVec
	registry *prometheus.Registry
}

func newRouter() *router {

	reg := prometheus.NewRegistry()

	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"handler", "method"},
	)

	reg.MustRegister(requestCounter)

	return &router{
		routes:   make(map[string]*route),
		metrics:  requestCounter,
		registry: reg,
	}
}

func InitRouter(handlers *handlers) http.Handler {
	router := newRouter()

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
	router.add("/metrics", http.MethodGet, router.handleMetrics)

	return router
}

func (router *router) add(path string, verbe string, callback func(http.ResponseWriter, *http.Request)) {
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
func (router *router) wrapHandlerWithMetrics(handlerName string, h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		router.metrics.WithLabelValues(handlerName, r.Method).Inc()
		h(w, r)
	}
}

// handleMetrics exposes Prometheus metrics
func (router *router) handleMetrics(w http.ResponseWriter, r *http.Request) {
	promhttp.HandlerFor(router.registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

func (router *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r, ok := router.routes[req.URL.Path]; ok {
		if handler, ok := r.handlers[req.Method]; ok {
			fmt.Println("ServeHTTP:", req.Method, req.URL.Path)
			handler(w, req)
			return
		}
		router.metrics.WithLabelValues(req.URL.Path, req.Method).Inc()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, req)
}
