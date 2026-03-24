package main

import (
	"net/http"

	"github.com/rs/zerolog/log"

	hellohandler "github.com/shiven-lohia/interneers-lab/pkg/helloworld/handler"

	productHandler "github.com/shiven-lohia/interneers-lab/pkg/products/handler"
	"github.com/shiven-lohia/interneers-lab/pkg/products/repository"
	productService "github.com/shiven-lohia/interneers-lab/pkg/products/service"

	"github.com/shiven-lohia/interneers-lab/pkg/middleware"
)

func main() {

	mux := http.NewServeMux()

	hellohandler.RegisterHelloHandler(mux)

	// PRODUCTS MODULE

	// repo := repository.NewMapProductRepository()
	client, err := repository.ConnectMongo()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	repo := repository.NewMongoProductRepository(client)

	pService := productService.NewProductService(repo)

	pHandler := productHandler.NewProductHandler(pService)

	productHandler.RegisterRoutes(mux, pHandler)

	// apply middleware
	loggedMux := middleware.LoggingMiddleware(mux)

	log.Info().Msg("Server starting on :8080")

	err = http.ListenAndServe(":8080", loggedMux)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
