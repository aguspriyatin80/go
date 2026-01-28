package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var products = []Product{
	{ID: 1, Name: "Fanta", Price: 5000, Stock: 50},
	{ID: 2, Name: "Sprite", Price: 5000, Stock: 60},
	{ID: 3, Name: "Coca-cola", Price: 5000, Stock: 70},
}

// DELETE  localhost:8080/api/products/{id} => API untuk update data berdasarkan ID
func deleteProductById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}
	for i, p := range products {
		if p.ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index
			products = append(products[:i], products[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}
	http.Error(w, "Invalid Request", http.StatusNotFound)
}

// PUT  localhost:8080/api/products/{id} => API untuk update data berdasarkan ID
func updateProductById(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}
	var updateProduct Product
	err = json.NewDecoder(r.Body).Decode(&updateProduct)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	for i := range products {
		if products[i].ID == id {
			updateProduct.ID = id
			products[i] = updateProduct
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduct)
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

// GET  localhost:8080/api/products/{id} => API untuk menampilkan satu data produk saja
func getProductById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid produk ID", http.StatusBadRequest) //404
		return
	}
	for _, p := range products {
		if p.ID == id {
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

// GET  localhost:8080/api/products => API untuk menampilkan semua data produk
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// POST localhost:8080/api/products => API untuk membuat produk baru
func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest) //404
	}
	newProduct.ID = len(products) + 1
	products = append(products, newProduct)
	w.WriteHeader(http.StatusCreated) //201
	json.NewEncoder(w).Encode(newProduct)
}

func main() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) { // pastikan endpoint tambahkan '/' setelah products
		if r.Method == http.MethodGet {
			getProductById(w, r)
		} else if r.Method == http.MethodPut {
			updateProductById(w, r)
		} else if r.Method == http.MethodDelete {
			deleteProductById(w, r)
		}
	})
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			getProducts(w, r)
		} else if r.Method == http.MethodPost {
			createProduct(w, r)
			return
		}
	})
	// localhost:8080/health => API untuk mengecek apakah server sudah running
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Server Running",
		})
	})
	fmt.Println("Server running di localhost: " + config.Port)
	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
