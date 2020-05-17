package controllers

import (
	"ben-jerry/models"
	productRepository "ben-jerry/repository/product"
	"ben-jerry/utils"
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

		products, _ = productRepo.GetProducts(db, product, products)
		utils.ResponseJSON(w, products)
	}
}

func (c Controller) GetProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		params := mux.Vars(r)

		productRepo := productRepository.ProductRepository{}

		id, err := strconv.Atoi(params["id"])
		if err != nil {
			var er models.Error
			er.Message = "ID required as int"
			utils.RespondWithError(w, http.StatusBadRequest, er)
			return
		}

		product, err = productRepo.GetProduct(db, product, id)
		if err != nil {
			var er models.Error
			er.Message = "No matching ID found "
			utils.RespondWithError(w, http.StatusBadRequest, er)
			return
		}
		utils.ResponseJSON(w, product)
	}
}

func (c Controller) AddProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		var productID int

		json.NewDecoder(r.Body).Decode(&product)

		productRepo := productRepository.ProductRepository{}
		productID = productRepo.AddProduct(db, product)

		utils.ResponseJSON(w, productID)
	}
}

func (c Controller) RemoveProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		productRepo := productRepository.ProductRepository{}

		id, err := strconv.Atoi(params["id"])
		if err != nil {
			var er models.Error
			er.Message = "ID required as int"
			utils.RespondWithError(w, http.StatusBadRequest, er)
			return
		}

		rowsDeleted := productRepo.RemoveProduct(db, id)

		json.NewEncoder(w).Encode(rowsDeleted)
	}
}

func (c Controller) UpdateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		params := mux.Vars(r)
		json.NewDecoder(r.Body).Decode(&product)

		id, err := strconv.Atoi(params["id"])
		if err != nil {
			var er models.Error
			er.Message = "ID required as int"
			utils.RespondWithError(w, http.StatusBadRequest, er)
			return
		}
		id_str := strconv.Itoa(id)
		productRepo := productRepository.ProductRepository{}
		rowsUpdated := productRepo.UpdateProduct(db, id_str, product)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}
func (c Controller) GetProductDetail(db *sql.DB, field int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		id, err := strconv.Atoi(params["id"])
		if err != nil {
			var er models.Error
			er.Message = "ID required as int"
			utils.RespondWithError(w, http.StatusBadRequest, er)
			return
		}
		productRepo := productRepository.ProductRepository{}
		value := productRepo.GetProductDetail(db, id, field)

		json.NewEncoder(w).Encode(value)
	}
}
func (c Controller) UpdateProductDetail(db *sql.DB, field int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			var er models.Error
			er.Message = "ID required as int"
			utils.RespondWithError(w, http.StatusBadRequest, er)
			return
		}

		var product models.Product
		err = json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			var er models.Error
			er.Message = "Parsing Error"
			utils.RespondWithError(w, http.StatusBadRequest, er)
			return
		}

		var newValue string

		if field == 0 {
			newValue = product.Name
		} else if field == 1 {
			newValue = product.ImageOpen
		} else if field == 2 {
			newValue = product.ImageClosed
		} else if field == 3 {
			newValue = product.Description
		} else if field == 4 {
			newValue = product.Story
		} else if field == 5 {
			newValue = product.AllergyInfo
		} else if field == 6 {
			newValue = product.DietaryCertifications
		} else {
			var er models.Error
			er.Message = "Field value exceeded"
			utils.RespondWithError(w, http.StatusBadRequest, er)
			return
		}

		log.Println(newValue)
		productRepo := productRepository.ProductRepository{}
		rowsUpdated := productRepo.UpdateProductDetail(db, id, field, newValue)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) GetProductIngredients(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		id := params["id"]

		productRepo := productRepository.ProductRepository{}
		ingredients := productRepo.GetProductIngredients(db, id)

		json.NewEncoder(w).Encode(ingredients)
	}
}

func (c Controller) UpdateProductIngredients(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["id"]

		var product models.Product
		_ = json.NewDecoder(r.Body).Decode(&product)

		newValue := product.Ingredients

		productRepo := productRepository.ProductRepository{}
		rowsUpdated := productRepo.UpdateProductIngredients(db, id, newValue)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) GetProductSourcingValues(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		id := params["id"]

		productRepo := productRepository.ProductRepository{}
		ingredients := productRepo.GetProductSourcingValues(db, id)

		json.NewEncoder(w).Encode(ingredients)
	}
}

func (c Controller) UpdateProductSourcingValues(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["id"]

		var product models.Product
		_ = json.NewDecoder(r.Body).Decode(&product)

		newValue := product.SourcingValues

		productRepo := productRepository.ProductRepository{}
		rowsUpdated := productRepo.UpdateProductSourcingValues(db, id, newValue)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

/*
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
}*/
