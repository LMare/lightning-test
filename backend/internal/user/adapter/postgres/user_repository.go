package postgres

import (
	"context"
	"database/sql"

	"github.com/Lmare/lightning-playground/backend/internal/user/domain"
	"github.com/Lmare/lightning-playground/backend/internal/user/port"
)

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

type PostgresUserRepository struct {
	db *sql.DB
}

var _ port.UserRepository = (*PostgresUserRepository)(nil)

func (r *PostgresUserRepository) FindAll(ctx context.Context) ([]domain.UserModel, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, nom, prenom, age, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.UserModel
	for rows.Next() {
		var p domain.UserModel
		if err := rows.Scan(&p.ID, &p.Nom, &p.Prenom, &p.Age, &p.Email); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	return list, nil
}
