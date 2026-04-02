package port

import (
	"context"

	"github.com/Lmare/lightning-playground/backend/internal/user/domain"
)

// UserRepository is the outbound port for user persistence.
type UserRepository interface {
	FindAll(ctx context.Context) ([]domain.UserModel, error)
}
