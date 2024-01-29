package main

import (
	"log"

	"github.com/luquxSentinel/checkout/service"
	"github.com/luquxSentinel/checkout/storage"
)

func main() {
	nopStorage := storage.NewNopStorage()
	productService := service.NewProductService(nopStorage)
	userService := service.NewUserService(nopStorage)

	api := NewAPIServer(":3000", productService, userService)

	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
