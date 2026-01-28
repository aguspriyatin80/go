package handlers

import (
	"belajar-go/services"
	"encoding/json"
	"net/http"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// HandleProducts - GET /api/produk
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	// switch r.Method {
	// case http.MethodGet:
	// 	h.GetAll(w, r)
	// case http.MethodPost:
	// 	h.Create(w, r)
	// default:
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// }
	if r.Method == http.MethodGet {
		h.GetAll(w, r)
	}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
