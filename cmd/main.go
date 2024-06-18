package main

import (
	"context"
	"log"
	"net/http"
    "shopapp/internal/products"
	"shopapp/config"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Database

func recoverHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("Recovered from panic: %v", rec)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}

func main() {
	client := config.ConnectToDB()
    defer func() {
        if err := client.Disconnect(context.TODO()); err != nil {
            log.Fatal(err)
        }
    }()

	collection := config.GetProductsCollection()
    productRepo := products.NewProductRepository(collection)
    productHandler := products.NewProductHandler(*productRepo)

	http.HandleFunc("/products/all", recoverHandler(productHandler.GetProductsHandler))
	http.HandleFunc("/products/add", recoverHandler(productHandler.AddProductHandler))
	
	log.Println("Starting server on :8020")
	err := http.ListenAndServe(":8020", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}