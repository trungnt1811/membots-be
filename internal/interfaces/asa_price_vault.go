package interfaces

import "context"

type TokenPriceRepo interface {
	GetAstraPrice(ctx context.Context) (int64, error)
}
