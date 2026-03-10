package main

import (
	"net/http"

	"github.com/rs/zerolog/log"

	hellohandler "github.com/shiven-lohia/interneers-lab/pkg/helloworld/handler"

	productController "github.com/shiven-lohia/interneers-lab/pkg/products/controller"
	productHandler "github.com/shiven-lohia/interneers-lab/pkg/products/handler"
	"github.com/shiven-lohia/interneers-lab/pkg/products/repository"

	"github.com/shiven-lohia/interneers-lab/pkg/middleware"
)

func main() {

	mux := http.NewServeMux()

	hellohandler.RegisterHelloHandler(mux)

	// PRODUCTS MODULE

	repo := repository.NewMapProductRepository()

	pController := productController.NewProductController(repo)

	pHandler := productHandler.NewProductHandler(pController)

	productHandler.RegisterRoutes(mux, pHandler)

	// apply middleware
	loggedMux := middleware.LoggingMiddleware(mux)

	log.Info().Msg("Server starting on :8080")

	err := http.ListenAndServe(":8080", loggedMux)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}