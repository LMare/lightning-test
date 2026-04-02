package personView

import "github.com/Lmare/lightning-playground/backend/internal/user/domain"

type PersonView struct {
	Index  int
	Person domain.UserModel
}

func ViewObject(persons []domain.UserModel) []PersonView {

	var viewData []PersonView
	for i, p := range persons {
		viewData = append(viewData, PersonView{
			Index:  i + 1,
			Person: p,
		})
	}
	return viewData

}
