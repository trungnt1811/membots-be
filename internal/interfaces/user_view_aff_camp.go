package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

type UserViewAffCampRepository interface {
	CreateUserViewAffCamp(ctx context.Context, data *model.UserViewAffCamp) error
}
