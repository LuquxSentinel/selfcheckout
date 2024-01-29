package service

import (
	"context"

	"github.com/luquxSentinel/checkout/types"
)

type ProductService interface {
	GetProductInfo(ctx context.Context, productID string) (*types.Product, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user *types.CreateUserInput) *types.Error
	Login(ctx context.Context, email string, password string) (*types.ResponseUser, *types.Error)
	AddProduct(ctx context.Context, productID string) *types.Error
	RemoveProduct(ctx context.Context, productID string) error
}
