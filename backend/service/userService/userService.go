package userService

import (
	"context"

	"github.com/Lmare/lightning-playground/backend/model/user"
	"github.com/Lmare/lightning-playground/backend/repository"
)

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type UserService struct {
	repo repository.UserRepository
}

func (s *UserService) List(ctx context.Context) ([]user.UserModel, error) {
	return s.repo.FindAll(ctx)
}

/*
func ListUsers() []user.UserModel {
	p1 := user.UserModel{Nom: "Dupont", Prenom: "Louis", Age: 29}
	p2 := user.New("Gedusor", "Tom", 21)
	p2.Prenom = "Voldemor"
	p3 := user.
		NewEmptyUserModel().
		SetNom("Soyer").
		SetPrenom("Tom").
		SetAge(8)

	var users []user.UserModel

	users = append(users, p1)
	users = append(users, p2)
	users = append(users, *p3)

	return users
}
*/
