package personView


import "github.com/Lmare/lightning-test/backend/model/personne"

type PersonView struct {
	Index int
	Person personne.Personne
}


func ViewObject(persons []personne.Personne) []PersonView {

	var viewData []PersonView
	for i, p := range persons {
	    viewData = append(viewData, PersonView{
	        Index:  i + 1,
	        Person: p,
	    })
	}
	return viewData

}
