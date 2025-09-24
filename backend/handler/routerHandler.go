package handler

import (
	"fmt"
	"net/http"
	"strings"
)

type route struct {
	path		string
	handlers 	map[string]func(http.ResponseWriter, *http.Request)
}

type Router struct {
	routes		map[string]*route
}

func GetRouter() *Router {
	router := Router{routes: make(map[string]*route)}
	router.add("/", http.MethodGet, HandleRoot)

	router.add("/lightning/alias", http.MethodPut, HandleUpdateNodeAlias)
	router.add("/lightning/channel", http.MethodPost, HandleOpenChannel)
	router.add("/lightning/nodes", http.MethodGet, HandleListOfNodes)
	router.add("/lightning/nodeInfo", http.MethodGet, HandleNodeInfo)
	router.add("/lightning/peer", http.MethodPost, HandleAddPeer)
	router.add("/lightning/uri", http.MethodGet, HandleShowUri)

	router.add("/users", http.MethodGet, HandleListPersonne)

	return &router
}


func (router *Router) add(path string, verbe string, callback func(http.ResponseWriter, *http.Request)) {
	if _, exist := router.routes[path]; !exist {
		router.routes[path] = &route{path:path, handlers: make(map[string]func(http.ResponseWriter, *http.Request))}
	}
	router.routes[path].handlers[verbe] = callback
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r, ok := router.routes[req.URL.Path]; ok {
		if handler, ok := r.handlers[req.Method]; ok {
			handler(w, req)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, req)
}



/**
 * Handler de / : propose une petite interface pour appeler les autres routes
 */
func HandleRoot(response http.ResponseWriter, request *http.Request) {

	html := `
	<h1>Lightning-test</h1>
	<p>
		Bienvenue sur le backend de cette application!
	</p>
	<p>
		Voici la liste des routes déservies :
	</p>
	<ul>`
	for _, r := range GetRouter().routes {
		html = html + "<li><a href=\"" + r.path + "\">" + r.path + "</a> ("
		// Extraire les méthodes
		methods := []string{}
		for method := range r.handlers {
			methods = append(methods, method)
		}
		html += strings.Join(methods, ", ") + ")</li>"
	}
	html = html + "</ul>"

	fmt.Fprintf(response, "%s", html)
}
