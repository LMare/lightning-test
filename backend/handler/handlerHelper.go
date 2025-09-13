package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"html/template"
	"errors"
    "log"
    "strings"
	"runtime"
)

// Transform the objet as a Json and put it in the reponse
func SetJsonResponse(response http.ResponseWriter, objet any) {

	objetJson, err := json.Marshal(objet)

	if err != nil {
		fmt.Println("Error json Marshal : ", err)
		fmt.Fprintf(response, "Une erreur est survenue")
	}
	response.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(response, "%s", string(objetJson))
}

// load the template HTML to render the object
func SetHtmlResponse(response http.ResponseWriter, file string, viewObject any) {
	response.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles(file))
	tmpl.Execute(response, viewObject)
}

// True if the client use HTMX, in this case the response muste be in HTML
func IsHTMX(request *http.Request) bool {
	return request.Header.Get("HX-Request") == "true"
}

// format the error for the logs
func LogException(err error) {
    if err == nil {
        return
    }

    // Capture les flags actuels
    originalFlags := log.Flags()

    // Première ligne : date, heure, fichier, ligne
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

    i := 1
    for err != nil {
        next := errors.Unwrap(err)

        currentMsg := err.Error()
        var nextMsg string
        if next != nil {
            nextMsg = next.Error()
        }

        context := strings.TrimSpace(strings.TrimSuffix(currentMsg, nextMsg))

        if i == 1 {
			_, file, line, _ := runtime.Caller(1)
			log.Printf("[ERROR] at %s:%d | %s", file, line, err.Error())
            log.SetFlags(0) // Supprime les métadonnées pour les lignes suivantes
        }

		log.Printf("\t→ Erreur #%d | Type: %T | Message : %s", i, err, context)

        err = next
        i++
    }

    // Restaure les flags initiaux
    log.SetFlags(originalFlags)
}
