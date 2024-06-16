package products

import (
	"encoding/json"
	"log"
	"net/http"
)

var productRepo *ProductRepository

func SetProductRepository(repo *ProductRepository) {
	productRepo = repo
}

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    products, err := productRepo.ListProducts()

    if err != nil {
		log.Fatal(err)
	}

    w.Header().Set("Content-Type", "application/json")
    
    if err := json.NewEncoder(w).Encode(products); err != nil {
        log.Printf("Failed to encode products to JSON: %v", err)
        http.Error(w, "Failed to encode products to JSON", http.StatusInternalServerError)
    }
}