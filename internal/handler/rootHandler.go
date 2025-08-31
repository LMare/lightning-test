package handler

import (
	"fmt"
	"net/http"
)

var initialised = false
var Routes []Route

type Route struct {
	Path     string
	Callback func(http.ResponseWriter, *http.Request)
}

func Init() {
	if !initialised {
		Routes = append(Routes, Route{Path: "/", Callback: HandleRoot})
		Routes = append(Routes, Route{Path: "/user", Callback: HandleListPersonne})

		initialised = true
	}
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
	for _, r := range Routes {
		html = html + "<li><a href=\"" + r.Path + "\">" + r.Path + "</a></li>"
	}
	html = html + "</ul>"

	fmt.Fprintf(response, "%s", html)
}
