package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SetJsonResponse(response http.ResponseWriter, objet any) {

	objetJson, err := json.Marshal(objet)

	if err != nil {
		fmt.Println("Error json Marshal : ", err)
		fmt.Fprintf(response, "Une erreur est survenue")
	}

	fmt.Fprintf(response, "%s", string(objetJson))
}
