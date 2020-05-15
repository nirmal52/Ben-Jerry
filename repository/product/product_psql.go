package productRepository

import (
	"ben-jerry/models"
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

type ProductRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getIngredients(db *sql.DB, productID string) []string {
	var ingredients []string
	rows1, err1 := db.Query("select value from ingredients where product_id=$1", productID)
	logFatal(err1)
	var ingredient string
	for rows1.Next() {
		err2 := rows1.Scan(&ingredient)
		logFatal(err2)
		ingredients = append(ingredients, ingredient)
	}
	return ingredients
}

func getSourcingValues(db *sql.DB, productID string) []string {
	var sourcingValues []string

	rows, err := db.Query("select value from sourcing_values where product_id=$1", productID)
	logFatal(err)
	var sourcingValue string
	for rows.Next() {
		err1 := rows.Scan(&sourcingValue)
		logFatal(err1)
		sourcingValues = append(sourcingValues, sourcingValue)
	}
	return sourcingValues
}

func (p ProductRepository) GetProducts(db *sql.DB, product models.Product, products []models.Product) []models.Product {
	rows, err := db.Query("select id, name, image_open, image_closed, description, story, allergy_info, dietary_certifications from product LIMIT 2")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&product.ProductId, &product.Name, &product.ImageOpen, &product.ImageClosed, &product.Description,
			&product.Story, &product.AllergyInfo, &product.DietaryCertifications)
		logFatal(err)

		product.Ingredients = getIngredients(db, product.ProductId)
		product.SourcingValues = getSourcingValues(db, product.ProductId)

		products = append(products, product)
	}

	return products
}

func (p ProductRepository) GetProduct(db *sql.DB, product models.Product, id int) models.Product {
	rows := db.QueryRow("select id, name, image_open, image_closed, description, story, allergy_info, dietary_certifications from product where id=$1", id)

	err := rows.Scan(&product.ProductId, &product.Name, &product.ImageOpen, &product.ImageClosed, &product.Description,
		&product.Story, &product.AllergyInfo, &product.DietaryCertifications)
	logFatal(err)

	product.Ingredients = getIngredients(db, product.ProductId)
	product.SourcingValues = getSourcingValues(db, product.ProductId)

	return product
}

func addIngredients(db *sql.DB, id int, ingredients []string) {
	var temp int
	for j := 0; j < len(ingredients); j++ {
		stmt := "insert into ingredients (product_id, value) values($1, $2) RETURNING i_id;"
		err := db.QueryRow(stmt, id, ingredients[j]).Scan(&temp)
		if err != nil {
			fmt.Println(err)
		}
	}
}
func addSourcingValues(db *sql.DB, id int, sourcing_values []string) {
	var temp int
	for j := 0; j < len(sourcing_values); j++ {
		stmt := "insert into sourcing_values (product_id, value) values($1, $2) RETURNING s_id;"
		err := db.QueryRow(stmt, id, sourcing_values[j]).Scan(&temp)
		if err != nil {
			fmt.Println(err)
		}
	}
}
func (p ProductRepository) AddProduct(db *sql.DB, product models.Product) int {
	var id int
	stmt := "insert into product (name, image_open, image_closed, description, story, allergy_info, dietary_certifications) values($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
	err := db.QueryRow(stmt, product.Name, product.ImageOpen, product.ImageClosed, product.Description, product.Story, product.AllergyInfo, product.DietaryCertifications).Scan(&id)
	addIngredients(db, id, product.Ingredients)
	addSourcingValues(db, id, product.SourcingValues)
	logFatal(err)

	return id
}

func deleteIngredients(db *sql.DB, id int) {
	_, _ = db.Exec("delete from ingredients where product_id = $1", id)
}
func deleteSourcingValues(db *sql.DB, id int) {
	_, _ = db.Exec("delete from sourcing_values where product_id = $1", id)
}

func (p ProductRepository) RemoveProduct(db *sql.DB, id int) int64 {
	deleteIngredients(db, id)
	deleteSourcingValues(db, id)
	result, err := db.Exec("delete from product where id = $1", id)
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	return rowsDeleted
}

func (p ProductRepository) UpdateProduct(db *sql.DB, product models.Product) int64 {
	id, _ := strconv.Atoi(product.ProductId)

	deleteIngredients(db, id)
	deleteSourcingValues(db, id)

	result, err := db.Exec("update product set name=$1,image_open =$2, image_closed=$3, description = $4, story = $5, allergy_info = $6, dietary_certifications = $7  where id=$8 RETURNING id",
		&product.Name, &product.ImageOpen, &product.ImageClosed, &product.Description, &product.Story, &product.AllergyInfo, &product.DietaryCertifications, &product.ProductId)

	addIngredients(db, id, product.Ingredients)
	addSourcingValues(db, id, product.SourcingValues)

	logFatal(err)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	return rowsUpdated
}

func (p ProductRepository) GetProductName(db *sql.DB, id int) string {
	var name string

	rows := db.QueryRow("select name from product where id=$1", id)

	err := rows.Scan(&name)
	logFatal(err)

	return name
}

func (p ProductRepository) UpdateProductName(db *sql.DB, id int, newName string) int64 {

	result, _ := db.Exec("update product set name=$1  where id=$2 RETURNING id",
		newName, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}

func (p ProductRepository) GetProductImageOpen(db *sql.DB, id int) string {
	var value string

	rows := db.QueryRow("select image_open from product where id=$1", id)

	err := rows.Scan(&value)
	logFatal(err)

	return value
}

func (p ProductRepository) UpdateProductImageOpen(db *sql.DB, id int, newValue string) int64 {

	result, _ := db.Exec("update product set image_open=$1  where id=$2 RETURNING id",
		newValue, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}

func (p ProductRepository) GetProductImageClose(db *sql.DB, id int) string {
	var value string

	rows := db.QueryRow("select image_closed from product where id=$1", id)

	err := rows.Scan(&value)
	logFatal(err)

	return value
}

func (p ProductRepository) UpdateProductImageClose(db *sql.DB, id int, newValue string) int64 {

	result, _ := db.Exec("update product set image_closed=$1  where id=$2 RETURNING id",
		newValue, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}

func (p ProductRepository) GetProductDescription(db *sql.DB, id int) string {
	var value string

	rows := db.QueryRow("select description from product where id=$1", id)

	err := rows.Scan(&value)
	logFatal(err)

	return value
}

func (p ProductRepository) UpdateProdutDdescription(db *sql.DB, id int, newValue string) int64 {

	result, _ := db.Exec("update product set description=$1  where id=$2 RETURNING id",
		newValue, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}

func (p ProductRepository) GetProductStory(db *sql.DB, id int) string {
	var value string

	rows := db.QueryRow("select story from product where id=$1", id)

	err := rows.Scan(&value)
	logFatal(err)

	return value
}

func (p ProductRepository) UpdateProductStory(db *sql.DB, id int, newValue string) int64 {

	result, _ := db.Exec("update product set story=$1  where id=$2 RETURNING id",
		newValue, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}

func (p ProductRepository) GetProductAllergy(db *sql.DB, id int) string {
	var value string

	rows := db.QueryRow("select allergy_info from product where id=$1", id)

	err := rows.Scan(&value)
	logFatal(err)

	return value
}

func (p ProductRepository) UpdateProductAllergy(db *sql.DB, id int, newValue string) int64 {

	result, _ := db.Exec("update product set allergy_info=$1  where id=$2 RETURNING id",
		newValue, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}

func (p ProductRepository) GetProductDiet(db *sql.DB, id int) string {
	var value string

	rows := db.QueryRow("select dietary_certifications from product where id=$1", id)

	err := rows.Scan(&value)
	logFatal(err)

	return value
}

func (p ProductRepository) UpdateProductDiet(db *sql.DB, id int, newValue string) int64 {

	result, _ := db.Exec("update product set dietary_certifications =$1  where id=$2 RETURNING id",
		newValue, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}
