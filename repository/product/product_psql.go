package productRepository

import (
	"ben-jerry/models"
	"database/sql"
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
func (p ProductRepository) GetProducts(db *sql.DB, product models.Product, products []models.Product) []models.Product {
	rows, err := db.Query("select id, name, image_open, image_closed, description, story, allergy_info, dietary_certifications from product LIMIT 5")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&product.ProductId, &product.Name, &product.ImageOpen, &product.ImageClosed, &product.Description,
			&product.Story, &product.AllergyInfo, &product.DietaryCertifications)
		logFatal(err)

		products = append(products, product)
	}

	return products
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
