package controllers

import (
	"ben-jerry/models"
	productRepository "ben-jerry/repository/product"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
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
