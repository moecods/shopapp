package products

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type ProductHandler struct {
    ProductRepo ProductRepository
}

func  NewProductHandler(repo ProductRepository) *ProductHandler {
    return &ProductHandler{ProductRepo: repo}
}


func (h *ProductHandler) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }


    products, err := h.ProductRepo.ListProducts()

    if err != nil {
		log.Fatal(err)
	}

    w.Header().Set("Content-Type", "application/json")
    
    if err := json.NewEncoder(w).Encode(products); err != nil {
        log.Printf("Failed to encode products to JSON: %v", err)
        http.Error(w, "Failed to encode products to JSON", http.StatusInternalServerError)
    }
}

func (h *ProductHandler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
    var product Product  

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Printf("Failed to read request body: %v", err)
        http.Error(w, "Failed to read request body", http.StatusInternalServerError)
        return
    }

    if err := json.Unmarshal(body, &product); err != nil {
        log.Printf("Failed to unmarshal request body: %v", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    err = h.ProductRepo.AddProduct(&product)
    if err != nil {
        http.Error(w, err.Error() , http.StatusBadRequest)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated) // 201 Created
    if err := json.NewEncoder(w).Encode(product); err != nil {
        log.Printf("Failed to encode response: %v", err)
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}