package productRepository

import (
	"ben-jerry/models"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type ProductRepository struct{}

func logFatal(err error, message string) {
	if err != nil {
		fmt.Println(message)
		fmt.Println(err)
	}
}

func getIngredients(db *sql.DB, productID string) []string {
	var ingredients []string
	//rows1, err1 := db.Query("select value from ingredients where product_id=$1", productID)
	rows1, err1 := db.Query("select value from ingredientsindex where id in (select value_id from ingredients where product_id = $1)", productID)
	if err1 != nil {
		fmt.Println("Error get ingredients")
		fmt.Println(productID)
		logFatal(err1, "get Ingredients 12")
	}
	var ingredient string
	for rows1.Next() {
		err2 := rows1.Scan(&ingredient)
		logFatal(err2, "get Ingredients 14")
		ingredients = append(ingredients, ingredient)
	}
	return ingredients
}

func getSourcingValues(db *sql.DB, productID string) []string {
	var sourcingValues []string

	//rows, err := db.Query("select value from sourcing_values where product_id=$1", productID)
	rows, err := db.Query("select value from sourcingvalueindex where id in (select value_id from sourcing_values where product_id = $1)", productID)
	logFatal(err, "get SV ")
	var sourcingValue string
	for rows.Next() {
		err1 := rows.Scan(&sourcingValue)
		logFatal(err1, "get SV ")
		sourcingValues = append(sourcingValues, sourcingValue)
	}
	return sourcingValues
}

func (p ProductRepository) GetProducts(db *sql.DB, product models.Product, products []models.Product) ([]models.Product, error) {
	rows, err := db.Query("select id, name, image_open, image_closed, description, story, allergy_info, dietary_certifications from products")

	if err != nil {
		if err != sql.ErrNoRows {
			return products, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&product.ProductId, &product.Name, &product.ImageOpen, &product.ImageClosed, &product.Description,
			&product.Story, &product.AllergyInfo, &product.DietaryCertifications)
		logFatal(err, "GetProducts SV ")

		product.Ingredients = getIngredients(db, product.ProductId)
		product.SourcingValues = getSourcingValues(db, product.ProductId)

		products = append(products, product)
	}

	return products, nil
}

func (p ProductRepository) GetProduct(db *sql.DB, product models.Product, id int) (models.Product, error) {
	rows := db.QueryRow("select id, name, image_open, image_closed, description, story, allergy_info, dietary_certifications from products where id=$1", id)

	err := rows.Scan(&product.ProductId, &product.Name, &product.ImageOpen, &product.ImageClosed, &product.Description,
		&product.Story, &product.AllergyInfo, &product.DietaryCertifications)

	if err != nil {
		if err == sql.ErrNoRows {
			return product, err
		}
	}

	product.Ingredients = getIngredients(db, product.ProductId)
	product.SourcingValues = getSourcingValues(db, product.ProductId)

	return product, nil
}

func getIngredientsIndex(db *sql.DB, values []string) []int {
	var newValues []int
	stmt := "insert into ingredientsindex(value) values($1) RETURNING id;"
	for i := 0; i < len(values); i++ {
		var newValue int
		row := db.QueryRow("select id from ingredientsindex where value = $1", values[i])
		err := row.Scan(&newValue)
		if err != nil {
			err2 := db.QueryRow(stmt, values[i]).Scan(&newValue)
			if err2 != nil {
				logFatal(err2, "Get ingredients index")
			} else {
				newValues = append(newValues, newValue)
			}
		} else {
			newValues = append(newValues, newValue)
		}
	}
	return newValues
}
func getSourceValueIndex(db *sql.DB, values []string) []int {
	var newValues []int
	stmt := "insert into sourcingvalueindex(value) values ( $1) RETURNING id;"
	for i := 0; i < len(values); i++ {
		var newValue int
		row := db.QueryRow("select id from sourcingvalueindex where value = $1", values[i])
		err := row.Scan(&newValue)
		if err != nil {
			err2 := db.QueryRow(stmt, values[i]).Scan(&newValue)
			if err2 != nil {
				logFatal(err, "getSourceValueIndex")
			} else {
				newValues = append(newValues, newValue)
			}
		} else {
			newValues = append(newValues, newValue)
		}
	}
	//fmt.Print("Returning values SV index search ")
	//fmt.Println(newValues)
	return newValues
}
func validateParenthesis(values []string) []string {
	var newValues []string
	for j := 0; j < len(values); j++ {
		if strings.Count(values[j], "(") != strings.Count(values[j], ")") {
			if j+1 != len(values) {
				values[j+1] = values[j] + values[j+1]
				continue
			}
		}
		newValues = append(newValues, values[j])
	}
	return newValues
}
func insertSourceValue(db *sql.DB, id string, values []string) {
	stmt := "insert into sourcing_values (product_id, value_id) values($1, $2) RETURNING id;"
	row_id := 0
	values = validateParenthesis(values)
	newValues := getSourceValueIndex(db, values)
	for j := 0; j < len(newValues); j++ {
		err := db.QueryRow(stmt, id, newValues[j]).Scan(&row_id)
		if err != nil {
			fmt.Print("ERRor 2 Source Values")
			logFatal(err, stmt)
		}
	}
}
func insertIngredients(db *sql.DB, id string, values []string) {
	row_id := 0
	stmt := "insert into ingredients (product_id, value_id) values($1, $2) RETURNING id;"
	values = validateParenthesis(values)
	newValues := getIngredientsIndex(db, values)
	//fmt.Println(newValues)
	for j := 0; j < len(newValues); j++ {
		err := db.QueryRow(stmt, id, newValues[j]).Scan(&row_id)
		if err != nil {
			fmt.Print("ERRor 1 ")
			fmt.Println(err)
		}
	}
}
func (p ProductRepository) AddProduct(db *sql.DB, product models.Product) int {
	var id int
	stmt := "insert into products (name, image_open, image_closed, description, story, allergy_info, dietary_certifications) values($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
	err := db.QueryRow(stmt, product.Name, product.ImageOpen, product.ImageClosed, product.Description, product.Story, product.AllergyInfo, product.DietaryCertifications).Scan(&id)
	if err != nil {
		logFatal(err, "Add Product")
	}
	id_str := strconv.Itoa(id)
	fmt.Println("new id  " + id_str)
	insertIngredients(db, id_str, product.Ingredients)
	insertSourceValue(db, id_str, product.SourcingValues)
	logFatal(err, "Add product ")

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
	result, err := db.Exec("delete from products where id = $1", id)
	logFatal(err, "REmove product")

	rowsDeleted, err := result.RowsAffected()
	logFatal(err, "remove product ")

	return rowsDeleted
}
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
func deleteIngredientsForUpdate(db *sql.DB, id string, valuesToDelete []string) {
	stmt := "delete from ingredients where product_id = $1 and value_id in ( select id from ingredientsindex where value = $2)"
	for i := 0; i < len(valuesToDelete); i++ {
		_, _ = db.Exec(stmt, id, valuesToDelete[i])
	}
}
func deleteSourceValuesForUpdate(db *sql.DB, id string, valuesToDelete []string) {
	stmt := "delete from sourcing_values where product_id = $1 and value_id in ( select id from sourcingvalueindex where value = $2)"
	for i := 0; i < len(valuesToDelete); i++ {
		_, _ = db.Exec(stmt, id, valuesToDelete[i])
	}
}
func updateIngredients(db *sql.DB, id string, newValues []string) {
	fmt.Println("UPDATE INGREDIENTS")
	oldValues := getIngredients(db, id)
	fmt.Println("Old Values")
	fmt.Println(oldValues)
	fmt.Println("new Values")
	fmt.Println(newValues)
	var toBeDeleted []string
	var toBeAdded []string
	for i := 0; i < len(oldValues); i++ {
		if !contains(newValues, oldValues[i]) {
			toBeDeleted = append(toBeDeleted, oldValues[i])
		}
	}
	for i := 0; i < len(newValues); i++ {
		if !contains(oldValues, newValues[i]) {
			toBeAdded = append(toBeAdded, newValues[i])
		}
	}
	fmt.Println(toBeDeleted)
	fmt.Println(toBeAdded)
	deleteIngredientsForUpdate(db, id, toBeDeleted)
	insertIngredients(db, id, toBeAdded)

}
func updateSourceValues(db *sql.DB, id string, newValues []string) {
	oldValues := getSourcingValues(db, id)
	fmt.Println("Old Values")
	fmt.Println(oldValues)
	fmt.Println("new Values")
	fmt.Println(newValues)
	var toBeDeleted []string
	var toBeAdded []string
	for i := 0; i < len(oldValues); i++ {
		if !contains(newValues, oldValues[i]) {
			toBeDeleted = append(toBeDeleted, oldValues[i])
		}
	}
	for i := 0; i < len(newValues); i++ {
		if !contains(oldValues, newValues[i]) {
			toBeAdded = append(toBeAdded, newValues[i])
		}
	}
	if len(oldValues) == 0 {
		toBeAdded = newValues
	}
	if len(newValues) == 0 {
		toBeDeleted = oldValues
	}
	fmt.Print("To be deleted   ")
	fmt.Println(toBeDeleted)
	fmt.Print("To be added    ")
	fmt.Println(toBeAdded)
	deleteSourceValuesForUpdate(db, id, toBeDeleted)
	insertSourceValue(db, id, toBeAdded)

}
func (p ProductRepository) UpdateProduct(db *sql.DB, productID string, product models.Product) int64 {
	//id, _ := strconv.Atoi(product.ProductId)

	//	deleteIngredients(db, id)
	//	deleteSourcingValues(db, id)
	fmt.Println("Product id")
	result, err := db.Exec("update products set name=$1,image_open =$2, image_closed=$3, description = $4, story = $5, allergy_info = $6, dietary_certifications = $7  where id=$8 RETURNING id",
		&product.Name, &product.ImageOpen, &product.ImageClosed, &product.Description, &product.Story, &product.AllergyInfo, &product.DietaryCertifications, &productID)

	if err != nil {
		return 0
	}

	updateIngredients(db, productID, product.Ingredients)
	updateSourceValues(db, productID, product.SourcingValues)
	//addIngredients(db, id, product.Ingredients)
	//addSourcingValues(db, id, product.SourcingValues)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err, "update product ")
	if err != nil {
		return 0
	}
	return rowsUpdated
}

func (p ProductRepository) GetProductDetail(db *sql.DB, id int, field int) string {
	identifyField := []string{"name", "image_open", "image_closed", "description", "story", "allergy_info", "dietary_certifications"}
	log.Println(len(identifyField))
	if field >= len(identifyField) {
		return ""
	}
	fieldToGet := identifyField[field]
	var value string
	rows := db.QueryRow("select "+fieldToGet+" from products where id=$1", id)

	err := rows.Scan(&value)

	if err != nil {
		return ""
	}

	return value
}
func (p ProductRepository) UpdateProductDetail(db *sql.DB, id int, field int, newName string) int64 {
	identifyField := []string{"name", "image_open", "image_closed", "description", "story", "allergy_info", "dietary_certifications"}
	log.Println(len(identifyField))
	if field >= len(identifyField) {
		return 0
	}
	fieldToGet := identifyField[field]

	result, _ := db.Exec("update products set "+fieldToGet+"=$1  where id=$2 RETURNING id",
		newName, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}
func (p ProductRepository) GetProductIngredients(db *sql.DB, id string) []string {
	ingredients := getIngredients(db, id)
	return ingredients
}

func (p ProductRepository) UpdateProductIngredients(db *sql.DB, id string, newValues []string) string {
	updateIngredients(db, id, newValues)
	return "1"
}

func (p ProductRepository) GetProductSourcingValues(db *sql.DB, id string) []string {
	sourcingValues := getSourcingValues(db, id)
	return sourcingValues
}

func (p ProductRepository) UpdateProductSourcingValues(db *sql.DB, id string, newValues []string) string {
	updateSourceValues(db, id, newValues)
	return "1"
}

/*
func (p ProductRepository) GetProductName(db *sql.DB, id int) string {
	var name string

	rows := db.QueryRow("select name from products where id=$1", id)

	err := rows.Scan(&name)
	logFatal(err, "get product name")

	return name
}

func (p ProductRepository) UpdateProductName(db *sql.DB, id int, newName string) int64 {

	result, _ := db.Exec("update products set name=$1  where id=$2 RETURNING id",
		newName, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}

func (p ProductRepository) GetProductImageOpen(db *sql.DB, id int) string {
	var value string

	rows := db.QueryRow("select image_open from products where id=$1", id)

	err := rows.Scan(&value)
	logFatal(err, "get image open")

	return value
}

func (p ProductRepository) UpdateProductImageOpen(db *sql.DB, id int, newValue string) int64 {

	result, _ := db.Exec("update products set image_open=$1  where id=$2 RETURNING id",
		newValue, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}

func (p ProductRepository) GetProductImageClose(db *sql.DB, id int) string {
	var value string

	rows := db.QueryRow("select image_closed from products where id=$1", id)

	err := rows.Scan(&value)
	logFatal(err, "get image close")

	return value
}

func (p ProductRepository) UpdateProductImageClose(db *sql.DB, id int, newValue string) int64 {

	result, _ := db.Exec("update productss set image_closed=$1  where id=$2 RETURNING id",
		newValue, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}

func (p ProductRepository) GetProductDescription(db *sql.DB, id int) string {
	var value string

	rows := db.QueryRow("select description from products where id=$1", id)

	err := rows.Scan(&value)
	logFatal(err, "get desvc")

	return value
}

func (p ProductRepository) UpdateProdutDdescription(db *sql.DB, id int, newValue string) int64 {

	result, _ := db.Exec("update productss set description=$1  where id=$2 RETURNING id",
		newValue, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}

func (p ProductRepository) GetProductStory(db *sql.DB, id int) string {
	var value string

	rows := db.QueryRow("select story from products where id=$1", id)

	err := rows.Scan(&value)
	logFatal(err, "story ")

	return value
}

func (p ProductRepository) UpdateProductStory(db *sql.DB, id int, newValue string) int64 {

	result, _ := db.Exec("update productss set story=$1  where id=$2 RETURNING id",
		newValue, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}

func (p ProductRepository) GetProductAllergy(db *sql.DB, id int) string {
	var value string

	rows := db.QueryRow("select allergy_info from products where id=$1", id)

	err := rows.Scan(&value)
	logFatal(err, "allergy ")

	return value
}

func (p ProductRepository) UpdateProductAllergy(db *sql.DB, id int, newValue string) int64 {

	result, _ := db.Exec("update productss set allergy_info=$1  where id=$2 RETURNING id",
		newValue, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}

func (p ProductRepository) GetProductDiet(db *sql.DB, id int) string {
	var value string

	rows := db.QueryRow("select dietary_certifications from products where id=$1", id)

	err := rows.Scan(&value)
	logFatal(err, "Diet ")

	return value
}

func (p ProductRepository) UpdateProductDiet(db *sql.DB, id int, newValue string) int64 {

	result, _ := db.Exec("update productss set dietary_certifications =$1  where id=$2 RETURNING id",
		newValue, id)

	rowsUpdated, _ := result.RowsAffected()

	return rowsUpdated
}
*/

/*
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
}*/
