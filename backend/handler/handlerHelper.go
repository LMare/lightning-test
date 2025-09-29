package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"html/template"
	"path/filepath"
	"errors"
    "log"
    "strings"
	"runtime"
)

func fail(response http.ResponseWriter, request *http.Request, message string, exception error) {
	LogException(exception)
	if IsHTMX(request) {
		HtmxMessageKo(response, message)
	} else {
		//TODO
	}
}

func htmxStreamEvent(response http.ResponseWriter, request *http.Request, streamId string) {
	if IsHTMX(request) {
		htmxResponse(response, "action/streamEvent.html", streamId)
	} else {
		jsonResponse(response, streamId)
	}
}


// Transform the objet as a Json and put it in the reponse
func jsonResponse(response http.ResponseWriter, objet any) {

	objetJson, err := json.Marshal(objet)

	if err != nil {
		fmt.Println("Error json Marshal : ", err)
		fmt.Fprintf(response, "Une erreur est survenue")
	}
	response.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(response, "%s", string(objetJson))
}

// render and object in Htmx
func htmxResponse(response http.ResponseWriter, file string, viewObject any) {
	response.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("backend/templates/" + file))
	tmpl.Execute(response, viewObject)
}

// render and object in Htmx with defined functions
func htmxResponseWithFuncs(response http.ResponseWriter, file string, viewObject any, funcs template.FuncMap) {
	response.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.New(filepath.Base(file)).
				Funcs(funcs).
				ParseFiles("backend/templates/" + file))
	tmpl.Execute(response, viewObject)
}


// return a response of success
// load the template HTML to render the object
func htmxMessageOk(response http.ResponseWriter, message string) {
	response.Header().Set("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("backend/templates/action/success.html"))
	tmpl.Execute(response, message)
}


// Retourne une réponse d'erreur
// ex : status = http.StatusForbidden
func HtmxMessageKo(response http.ResponseWriter, message string ) {
	response.Header().Set("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("backend/templates/action/error.html"))
	tmpl.Execute(response, message)
}

// return code 204
func OkNoContent(response http.ResponseWriter){
	response.WriteHeader(http.StatusNoContent)
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

        context := strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(currentMsg, nextMsg), "→"))

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
