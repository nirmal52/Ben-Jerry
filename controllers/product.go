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

func (c Controller) UpdateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		json.NewDecoder(r.Body).Decode(&product)

		productRepo := productRepository.ProductRepository{}
		rowsUpdated := productRepo.UpdateProduct(db, product)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) GetProductName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		productRepo := productRepository.ProductRepository{}
		title := productRepo.GetProductName(db, id)

		json.NewEncoder(w).Encode(title)
	}
}

func (c Controller) UpdateProductName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])

		var product models.Product
		_ = json.NewDecoder(r.Body).Decode(&product)

		newName := product.Name

		productRepo := productRepository.ProductRepository{}
		rowsUpdated := productRepo.UpdateProductName(db, id, newName)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) GetProductImageOpen(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		productRepo := productRepository.ProductRepository{}
		title := productRepo.GetProductImageOpen(db, id)

		json.NewEncoder(w).Encode(title)
	}
}

func (c Controller) UpdateProductImageOpen(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])

		var product models.Product
		_ = json.NewDecoder(r.Body).Decode(&product)

		newValue := product.ImageOpen

		productRepo := productRepository.ProductRepository{}
		rowsUpdated := productRepo.UpdateProductImageOpen(db, id, newValue)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) GetProductImageClose(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		productRepo := productRepository.ProductRepository{}
		title := productRepo.GetProductImageClose(db, id)

		json.NewEncoder(w).Encode(title)
	}
}

func (c Controller) UpdateProductImageClose(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])

		var product models.Product
		_ = json.NewDecoder(r.Body).Decode(&product)

		newValue := product.ImageClosed

		productRepo := productRepository.ProductRepository{}
		rowsUpdated := productRepo.UpdateProductImageClose(db, id, newValue)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) GetProductDescription(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		productRepo := productRepository.ProductRepository{}
		title := productRepo.GetProductDescription(db, id)

		json.NewEncoder(w).Encode(title)
	}
}

func (c Controller) UpdateProdutDdescription(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])

		var product models.Product
		_ = json.NewDecoder(r.Body).Decode(&product)

		newValue := product.Description

		productRepo := productRepository.ProductRepository{}
		rowsUpdated := productRepo.UpdateProdutDdescription(db, id, newValue)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) GetProductStory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		productRepo := productRepository.ProductRepository{}
		title := productRepo.GetProductStory(db, id)

		json.NewEncoder(w).Encode(title)
	}
}

func (c Controller) UpdateProductStory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])

		var product models.Product
		_ = json.NewDecoder(r.Body).Decode(&product)

		newValue := product.Story

		productRepo := productRepository.ProductRepository{}
		rowsUpdated := productRepo.UpdateProductStory(db, id, newValue)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) GetProductAllergy(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		productRepo := productRepository.ProductRepository{}
		title := productRepo.GetProductAllergy(db, id)

		json.NewEncoder(w).Encode(title)
	}
}

func (c Controller) UpdateProductAllergy(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])

		var product models.Product
		_ = json.NewDecoder(r.Body).Decode(&product)

		newValue := product.AllergyInfo

		productRepo := productRepository.ProductRepository{}
		rowsUpdated := productRepo.UpdateProductAllergy(db, id, newValue)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) GetProductDiet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		productRepo := productRepository.ProductRepository{}
		title := productRepo.GetProductDiet(db, id)

		json.NewEncoder(w).Encode(title)
	}
}

func (c Controller) UpdateProductDiet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])

		var product models.Product
		_ = json.NewDecoder(r.Body).Decode(&product)

		newValue := product.DietaryCertifications

		productRepo := productRepository.ProductRepository{}
		rowsUpdated := productRepo.UpdateProductDiet(db, id, newValue)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}
