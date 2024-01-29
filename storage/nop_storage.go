package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/luquxSentinel/checkout/types"
)

type NopStorage struct {
	products []*types.Product
	users    []*types.User
}

func NewNopStorage() *NopStorage {
	products := make([]*types.Product, 0)
	products = append(products, &types.Product{
		ID:          fmt.Sprintf("%d", len(products)+1),
		BrandName:   "Nestle",
		ProductName: "Barone",
		ExpiryDate:  time.Now().AddDate(2, 1, 3),
		Price:       11.00,
	})

	return &NopStorage{
		products: products,
		users:    make([]*types.User, 0),
	}
}

func (s *NopStorage) GetProductInfo(ctx context.Context, productID string) (*types.Product, error) {
	for _, product := range s.products {
		if product.ID == productID {
			return product, nil
		}
	}

	return nil, fmt.Errorf("product with ID:%s not found", productID)
}

func (s *NopStorage) GetBusket(ctx context.Context, uid string) (*types.Busket, error) {
	for _, user := range s.users {
		if user.UID == uid {
			return user.Busket, nil
		}
	}

	return nil, fmt.Errorf("busket not found")
}

func (s *NopStorage) UpdateBusket(ctx context.Context, uid string, newBusket *types.Busket) error {
	for idx, user := range s.users {
		if user.UID == uid {
			s.users[idx].Busket = newBusket
			return nil
		}
	}

	return fmt.Errorf("user not found")
}

func (s *NopStorage) CheckEmailExists(ctx context.Context, email string) (int64, error) {

	for _, user := range s.users {
		if user.Email == email {
			return 1, nil
		}
	}
	return 0, nil
}

func (s *NopStorage) CreateUser(ctx context.Context, user *types.User) error {
	s.users = append(s.users, user)
	return nil
}

func (s *NopStorage) GetUser(ctx context.Context, email string) (*types.User, error) {
	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}
