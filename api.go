package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/luquxSentinel/checkout/service"
	"github.com/luquxSentinel/checkout/types"
)

type APIFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) *types.Error
type APIServer struct {
	productService service.ProductService
	userService    service.UserService
	listenAddress  string
	router         *mux.Router
}

func NewAPIServer(listenAddress string, productService service.ProductService, userService service.UserService) *APIServer {
	return &APIServer{
		userService:    userService,
		productService: productService,
		listenAddress:  listenAddress,
		router:         mux.NewRouter(),
	}
}

func (api *APIServer) Run() error {

	api.router.HandleFunc("/get-product/{product-id}", api.handlerFunc(api.getProductInfo))

	fmt.Printf("Listening on port : %s", strings.Split(api.listenAddress, ":")[1])
	return http.ListenAndServe(api.listenAddress, api.router)
}

func (api *APIServer) handlerFunc(fn APIFunc) http.HandlerFunc {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request-id", uuid.New())
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(ctx, w, r)
		if err != nil {
			api.WriteError(w, err)
		}
	}
}

func (api *APIServer) login(ctx context.Context, w http.ResponseWriter, r *http.Request) *types.Error {
	loginInput := new(types.LoginInput)

	if err := BodyParse(r.Body, loginInput); err != nil {
		return &types.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid login body",
		}
	}

	user, err := api.userService.Login(ctx, loginInput.Email, loginInput.Password)
	if err != nil {
		return err
	}

	api.WriteSuccess(w, user)
	return nil
}

func (api *APIServer) createUser(ctx context.Context, w http.ResponseWriter, r *http.Request) *types.Error {
	createUserInput := new(types.CreateUserInput)
	if err := BodyParse(r.Body, createUserInput); err != nil {
		return &types.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid create user data",
		}
	}

	err := api.userService.CreateUser(ctx, createUserInput)
	if err != nil {
		return err
	}

	api.WriteSuccess(w, map[string]string{"message": "user successfully created"})
	return nil
}

func (api *APIServer) getProductInfo(ctx context.Context, w http.ResponseWriter, r *http.Request) *types.Error {
	productID, ok := mux.Vars(r)["product-id"]
	if !ok {
		return &types.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid product id",
		}
	}

	product, err := api.productService.GetProductInfo(ctx, productID)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	// TODO:handle write error
	api.WriteSuccess(w, product)
	return nil
}

func (api *APIServer) addToBusket(ctx context.Context, w http.ResponseWriter, r *http.Request) *types.Error {
	productID, ok := mux.Vars(r)["product-id"]
	if !ok {
		return &types.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid product id",
		}
	}

	err := api.userService.AddProduct(ctx, productID)
	if err != nil {
		return err
	}

	api.WriteSuccess(w, map[string]string{"message": "successfully added product to busket"})
	return nil
}

func (api *APIServer) WriteError(w http.ResponseWriter, err *types.Error) {
	w.WriteHeader(err.StatusCode)
	b, _ := json.Marshal(map[string]string{"error": err.Error()})
	fmt.Fprintf(w, "%v", b)
}

func (api *APIServer) WriteSuccess(w http.ResponseWriter, v any) {
	w.WriteHeader(http.StatusOK)
	// b, _ := json.Marshal(product)
	// fmt.Fprintf(w, "%v", json.Encoder().)
	json.NewEncoder(w).Encode(v)
}

func BodyParse(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
