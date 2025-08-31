
package main


import "fmt"
import "github.com/Lmare/lightning-test/internal/model/personne"


func main() {
	fmt.Println("Hello World")

	//p1 := personne.Personne{Nom : "Dupont", Prenom : "Louis", Age : 29,}
	p2 := personne.New("Gedusor", "Tom", 21)
	//p2.Prenom = "Voldemor"
	p3 := personne.
			NewEmptyPersonne().
			SetNom("Soyer").
			SetPrenom("Tom").
			SetAge(8);
	//fmt.Println(p1)
	fmt.Println(p2)
	fmt.Println(p3)
	fmt.Println(*p3)
}