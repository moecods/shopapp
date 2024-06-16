package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"net/http"
    "shopapp/internal/products"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

func connectToDB() {
	uri := "mongodb://localhost:27017"

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

	err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
		os.Exit(1)
    }

	database := client.Database("shop_db")
    collection = database.Collection("products")
	fmt.Println("Connected to MongoDB!")
}

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
	connectToDB()
	
	http.HandleFunc("/products/all", recoverHandler(products.GetProductsHandler))
	
	log.Println("Starting server on :8020")
	err := http.ListenAndServe(":8020", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
	defer client.Disconnect(context.TODO())
}