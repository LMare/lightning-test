package service

import "github.com/Lmare/lightning-test/backend/model/personne"

func ListUsers() []personne.Personne {
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

	return users
}
