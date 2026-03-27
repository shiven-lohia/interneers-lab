package main

import (
	"net/http"

	"github.com/rs/zerolog/log"

	hellohandler "github.com/shiven-lohia/interneers-lab/pkg/helloworld/handler"
	"github.com/shiven-lohia/interneers-lab/pkg/middleware"
	productHandler "github.com/shiven-lohia/interneers-lab/pkg/products/handler"
	"github.com/shiven-lohia/interneers-lab/pkg/products/repository"
	productService "github.com/shiven-lohia/interneers-lab/pkg/products/service"
)

func main() {

	mux := http.NewServeMux()

	hellohandler.RegisterHelloHandler(mux)

	// PRODUCTS MODULE

	client, err := repository.ConnectMongo()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	// product repo
	pRepo := repository.NewMongoProductRepository(client)

	// category repo
	cCollection := client.Database("inventory").Collection("categories")
	cRepo := repository.NewMongoCategoryRepository(cCollection)

	// services
	cService := productService.NewProductCategoryService(cRepo)
	pService := productService.NewProductService(pRepo, cRepo)

	// handlers
	cHandler := productHandler.NewProductCategoryHandler(cService)
	pHandler := productHandler.NewProductHandler(pService)

	// routes
	productHandler.RegisterCategoryRoutes(mux, cHandler)
	productHandler.RegisterRoutes(mux, pHandler)

	// apply middleware
	loggedMux := middleware.LoggingMiddleware(mux)
	corsMux := middleware.CORSMiddleware(loggedMux)

	log.Info().Msg("Server starting on :8080")

	err = http.ListenAndServe(":8080", corsMux)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
