package repository

import (
	"context"
	"database/sql"

	"github.com/Lmare/lightning-playground/backend/model/user"
)

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

type UserRepository interface {
	FindAll(ctx context.Context) ([]user.UserModel, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func (r *PostgresUserRepository) FindAll(ctx context.Context) ([]user.UserModel, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, nom, prenom, age, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []user.UserModel
	for rows.Next() {
		var p user.UserModel
		if err := rows.Scan(&p.ID, &p.Nom, &p.Prenom, &p.Age, &p.Email); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	return list, nil
}
