package controllers

import (
	"ben-jerry/models"
	productRepository "ben-jerry/repository/product"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var products []models.Product

type Controller struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c Controller) GetProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		products = []models.Product{}
		productRepo := productRepository.ProductRepository{}

		products = productRepo.GetProducts(db, product, products)

		json.NewEncoder(w).Encode(products)
	}
}

func (c Controller) GetProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		params := mux.Vars(r)

		productRepo := productRepository.ProductRepository{}

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		product = productRepo.GetProduct(db, product, id)

		json.NewEncoder(w).Encode(product)
	}
}

func (c Controller) AddProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		var productID int

		json.NewDecoder(r.Body).Decode(&product)

		productRepo := productRepository.ProductRepository{}
		productID = productRepo.AddProduct(db, product)

		json.NewEncoder(w).Encode(productID)
	}
}

func (c Controller) RemoveProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		productRepo := productRepository.ProductRepository{}

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		rowsDeleted := productRepo.RemoveProduct(db, id)

		json.NewEncoder(w).Encode(rowsDeleted)
	}
}
