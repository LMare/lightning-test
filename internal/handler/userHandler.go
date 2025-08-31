package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Lmare/lightning-test/internal/model/personne"
)

func HandleListPersonne(response http.ResponseWriter, request *http.Request) {
	p1 := personne.Personne{Nom: "Dupont", Prenom: "Louis", Age: 29}
	p2 := personne.New("Gedusor", "Tom", 21)
	p2.Prenom = "Voldemor"
	p3 := personne.
		NewEmptyPersonne().
		SetNom("Soyer").
		SetPrenom("Tom").
		SetAge(8)

	var users []personne.Personne

	users = append(users, p1)
	users = append(users, p2)
	users = append(users, *p3)

	fmt.Println(users)
	usersJson, err := json.Marshal(users)

	if err != nil {
		fmt.Println("Error json Marshal : ", err)
		fmt.Fprintf(response, "Une erreur est survenue")
	}

	fmt.Fprintf(response, string(usersJson))
}
