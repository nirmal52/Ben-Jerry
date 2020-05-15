package main

import (
	"ben-jerry/controllers"
	"ben-jerry/driver"
	"ben-jerry/models"
	"ben-jerry/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var db *sql.DB

func init() {
	gotenv.Load()
}

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func uploadJSONdataToDB() {
	file, _ := ioutil.ReadFile("icecream.json")
	var newProduct []models.Product
	_ = json.Unmarshal([]byte(file), &newProduct)
	for i := 0; i < len(newProduct); i++ {
		id := 0
		curr := newProduct[i]
		stmt := "insert into product (id, name, image_open, image_closed, description, story, allergy_info, dietary_certifications) values($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;"
		err := db.QueryRow(stmt, curr.ProductId, curr.Name, curr.ImageOpen, curr.ImageClosed, curr.Description, curr.Story, curr.AllergyInfo, curr.DietaryCertifications).Scan(&id)
		if err != nil {
			fmt.Println(err)
			fmt.Println(id)
		}
		for j := 0; j < len(curr.Ingredients); j++ {
			stmt := "insert into ingredients (product_id, value) values($1, $2) RETURNING i_id;"
			err := db.QueryRow(stmt, curr.ProductId, curr.Ingredients[j]).Scan(&id)
			if err != nil {
				fmt.Println(err)
			}
		}
		for j := 0; j < len(curr.SourcingValues); j++ {
			stmt := "insert into sourcing_values (product_id, value) values($1, $2) RETURNING s_id;"
			err := db.QueryRow(stmt, curr.ProductId, curr.SourcingValues[j]).Scan(&id)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}

func main() {
	db = driver.ConnectDB()
	router := mux.NewRouter()

	controller := controllers.Controller{}

	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")

	router.HandleFunc("/books", controller.GetProducts(db)).Methods("GET")

	router.HandleFunc("/protectedEndpoint", utils.TokenVerifyMiddleWare(controller.Protected(db))).Methods("GET")

	router.HandleFunc("/products", utils.TokenVerifyMiddleWare(controller.GetProducts(db))).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
	//uploadJSONdataToDB()
}
