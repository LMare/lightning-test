package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"html/template"
)

// Transform l'objet en Json et le met dans la reponse
func SetJsonResponse(response http.ResponseWriter, objet any) {

	objetJson, err := json.Marshal(objet)

	if err != nil {
		fmt.Println("Error json Marshal : ", err)
		fmt.Fprintf(response, "Une erreur est survenue")
	}
	response.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(response, "%s", string(objetJson))
}

func SetHtmlResponse(response http.ResponseWriter, file string, viewObject any) {
	response.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles(file))
	tmpl.Execute(response, viewObject)
}


func IsHTMX(request *http.Request) bool {
	return request.Header.Get("HX-Request") == "true"
}
