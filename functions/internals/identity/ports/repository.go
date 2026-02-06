package ports

import (
	"context"

	"github.com/a-ioenimil/serverless-dispatch/functions/internals/identity/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user domain.User) error
	Get(ctx context.Context, id string) (*domain.User, error)
}
