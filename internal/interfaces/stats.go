package interfaces

import (
	"context"
)

type StatsUCase interface {
	GetStatsByMemeAddress(ctx context.Context, address string) (interface{}, error)
}
