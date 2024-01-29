package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/luquxSentinel/checkout/storage"
	"github.com/luquxSentinel/checkout/types"
)

type UserServiceImpl struct {
	storage storage.Storage
}

func NewUserService(storage storage.Storage) *UserServiceImpl {
	return &UserServiceImpl{
		storage: storage,
	}
}

func (s *UserServiceImpl) Login(ctx context.Context, email string, password string) (*types.ResponseUser, *types.Error) {
	user, err := s.storage.GetUser(ctx, email)
	if err != nil {
		return nil, &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	err = verifyPassword(user.Password, password)
	if err != nil {
		return nil, &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	return user.ResponseUser(), nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *types.CreateUserInput) *types.Error {
	emailCount, err := s.storage.CheckEmailExists(ctx, user.Email)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "error validating email",
		}
	}

	if emailCount > 0 {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "email already in use",
		}
	}

	newUser := new(types.User)
	newUser.UID = uuid.New().String()
	newUser.Email = user.Email
	newUser.FirstName = user.FirstName
	newUser.LastName = user.LastName

	// TODO: hash password
	newUser.Password = user.Password
	newUser.Busket = types.NewBusket(newUser.UID)
	newUser.CreatedAt = time.Now().UTC()

	err = s.storage.CreateUser(ctx, newUser)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "error creating user",
		}
	}
	return nil
}

func (s *UserServiceImpl) AddProduct(ctx context.Context, uid, productID string) *types.Error {

	busket, err := s.storage.GetBusket(ctx, uid)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	for idx, product := range busket.Products {
		if product.ID == productID {
			busket.Products[idx].UnitSize += 1
			err = s.storage.UpdateBusket(ctx, uid, busket)
			if err != nil {
				return &types.Error{
					StatusCode: http.StatusInternalServerError,
					Message:    err.Error(),
				}
			}
			return nil
		}
	}

	product, err := s.storage.GetProductInfo(ctx, productID)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	newBusketProduct := &types.BusketProduct{
		ID:        product.ID,
		UnitPrice: product.Price,
		UnitSize:  1,
	}

	busket.Products = append(busket.Products, newBusketProduct)

	err = s.storage.UpdateBusket(ctx, uid, busket)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	return nil
}

func (s *UserServiceImpl) RemoveProduct(ctx context.Context, productID string) error {
	return nil
}

func (s *UserServiceImpl) GetBusket(ctx context.Context, uid string) (*types.Busket, *types.Error) {
	busket, err := s.storage.GetBusket(ctx, uid)
	if err != nil {
		return nil, &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}

	}
	return busket, nil
}

func verifyPassword(foundPassword, givenPassword string) error {
	if foundPassword == givenPassword {
		return nil
	}

	return fmt.Errorf("wrong email or password")
}
