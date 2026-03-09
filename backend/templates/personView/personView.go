package personView

import "github.com/Lmare/lightning-playground/backend/model/user"

type PersonView struct {
	Index  int
	Person user.UserModel
}

func ViewObject(persons []user.UserModel) []PersonView {

	var viewData []PersonView
	for i, p := range persons {
		viewData = append(viewData, PersonView{
			Index:  i + 1,
			Person: p,
		})
	}
	return viewData

}
