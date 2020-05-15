package productRepository

import (
	"ben-jerry/models"
	"database/sql"
	"fmt"
	"log"
)

type ProductRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/*
create table product (
  id serial primary key,
  name text,
  image_open text,
  image_closed text,
  description text,
  story text,
  allergy_info text,
  dietary_certifications text
);


*/

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

func (p ProductRepository) RemoveProduct(db *sql.DB, id int) int64 {
	_, _ = db.Exec("delete from ingredients where product_id = $1", id)

	_, _ = db.Exec("delete from sourcing_values where product_id = $1", id)

	result, err := db.Exec("delete from product where id = $1", id)
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	return rowsDeleted
}

/*
Comment
*/
/*func (b BookRepository) GetBook(db *sql.DB, book models.Book, id int) models.Book {
	rows := db.QueryRow("select * from books where id=$1", id)

	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)

	return book
}

func (b BookRepository) AddBook(db *sql.DB, book models.Book) int {
	err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;",
		book.Title, book.Author, book.Year).Scan(&book.ID)

	logFatal(err)

	return book.ID
}

func (b BookRepository) UpdateBook(db *sql.DB, book models.Book) int64 {
	result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id",
		&book.Title, &book.Author, &book.Year, &book.ID)

	logFatal(err)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	return rowsUpdated
}

func (b BookRepository) RemoveBook(db *sql.DB, id int) int64 {
	result, err := db.Exec("delete from books where id = $1", id)
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	return rowsDeleted
}
*/

//rows1, err1 := db.Query("select value from ingredients where product_id=$1", product.ProductId)
//logFatal(err1)
//product.Ingredients = getIngredients(db, product.ProductId)
//product.SourcingValues = getSourcingValues(db, product.ProductId)
/*var ingredient string
for rows1.Next() {
	err2 := rows1.Scan(&ingredient)
	logFatal(err2)
	product.Ingredients = append(product.Ingredients, ingredient)
}

rows2, err3 := db.Query("select value from sourcing_values where product_id=$1", product.ProductId)
logFatal(err3)
var sourcing_value string
for rows2.Next() {
	err4 := rows2.Scan(&sourcing_value)
	logFatal(err4)
	product.SourcingValues = append(product.SourcingValues, sourcing_value)
}
*/
