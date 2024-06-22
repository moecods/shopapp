package products

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

// UpdateProductHandler handles updating an existing product.
func (h *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var product Product
    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    id, err := primitive.ObjectIDFromHex(product.ID.Hex())
 
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    if err := h.ProductRepo.UpdateProduct(id, &product); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(product)
}