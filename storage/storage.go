package storage

import (
	"context"

	"github.com/luquxSentinel/checkout/types"
)

type Storage interface {
	GetProductInfo(ctx context.Context, productID string) (*types.Product, error)
	GetBusket(ctx context.Context, uid string) (*types.Busket, error)
	UpdateBusket(ctx context.Context, uid string, newBusket *types.Busket) error
	CheckEmailExists(ctx context.Context, email string) (int64, error)
	CreateUser(ctx context.Context, user *types.User) error
	GetUser(ctx context.Context, email string) (*types.User, error)
	AddProduct(ctx context.Context, productID string) error
}
