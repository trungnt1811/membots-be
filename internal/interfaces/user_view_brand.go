package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

type UserViewBrandRepository interface {
	CreateUserViewBrand(ctx context.Context, data *model.UserViewBrand) error
}
