package service

import (
	"context"

	"github.com/luquxSentinel/checkout/storage"
	"github.com/luquxSentinel/checkout/types"
)

type ProductServiceImpl struct {
	storage storage.Storage
}

func NewProductService(storage storage.Storage) *ProductServiceImpl {
	return &ProductServiceImpl{
		storage: storage,
	}
}

func (s ProductServiceImpl) GetProductInfo(ctx context.Context, productID string) (*types.Product, error) {
	return s.storage.GetProductInfo(ctx, productID)
}
