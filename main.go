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
	"strings"

	"ben-jerry/docs"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Ben - Jerry API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationUrl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationUrl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

var db *sql.DB

func init() {
	gotenv.Load()
}

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getIngredientsIndex(values []string) []int {
	var newValues []int
	stmt := "insert into ingredientsindex(value) values($1) RETURNING id;"
	for i := 0; i < len(values); i++ {
		var newValue int
		row := db.QueryRow("select id from ingredientsindex where value = $1", values[i])
		err := row.Scan(&newValue)
		if err != nil {
			logError(err)
			//		fmt.Println("Trying to Insert")
			err2 := db.QueryRow(stmt, values[i]).Scan(&newValue)
			if err2 != nil {
				//			fmt.Println("Error 11")
				logError(err2)
			} else {
				//fmt.Println(newValue)
				newValues = append(newValues, newValue)
			}
		} else {
			newValues = append(newValues, newValue)
		}
	}
	return newValues
}
func getSourceValueIndex(values []string) []int {
	var newValues []int
	stmt := "insert into sourcingvalueindex(value) values ( $1) RETURNING id;"
	//fmt.Println("Inside Sourcing values ")
	for i := 0; i < len(values); i++ {
		//fmt.Println(values[i])
		var newValue int
		row := db.QueryRow("select id from sourcingvalueindex where value = $1", values[i])
		err := row.Scan(&newValue)
		if err != nil {
			logError(err)
			//		fmt.Println("Trying to Insert SV")
			err2 := db.QueryRow(stmt, values[i]).Scan(&newValue)
			if err2 != nil {
				//			fmt.Println("Error 11 SV ")
				logError(err2)
			} else {
				//			fmt.Println(newValue)
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
func insertIntoDB(products []models.Product) {
	for i := 0; i < len(products); i++ {
		insertStringValuesToDB(products[i])
		insertIngredients(products[i].ProductId, products[i].Ingredients)
		insertSourceValue(products[i].ProductId, products[i].SourcingValues)
	}
}
func insertStringValuesToDB(product models.Product) {
	id := 0
	stmt := "insert into products (id, name, image_open, image_closed, description, story, allergy_info, dietary_certifications) values($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;"
	err := db.QueryRow(stmt, product.ProductId, product.Name, product.ImageOpen, product.ImageClosed, product.Description, product.Story, product.AllergyInfo, product.DietaryCertifications).Scan(&id)
	if err != nil {
		fmt.Print("ERRor 3 ")
		fmt.Println(err)
	}
}
func insertSourceValue(id string, values []string) {
	stmt := "insert into sourcing_values (product_id, value_id) values($1, $2) RETURNING id;"
	row_id := 0
	values = validateParenthesis(values)
	//fmt.Print("Validate parenthses")
	//fmt.Println(values)
	newValues := getSourceValueIndex(values)
	//fmt.Print("SOurcing VAlues  ")
	//fmt.Println(newValues)
	for j := 0; j < len(newValues); j++ {
		err := db.QueryRow(stmt, id, newValues[j]).Scan(&row_id)
		if err != nil {
			fmt.Print("ERRor 2 ")
			logError(err)
		}
	}
}
func insertIngredients(id string, values []string) {
	row_id := 0
	stmt := "insert into ingredients (product_id, value_id) values($1, $2) RETURNING id;"
	values = validateParenthesis(values)
	newValues := getIngredientsIndex(values)
	//fmt.Println(newValues)
	for j := 0; j < len(newValues); j++ {
		err := db.QueryRow(stmt, id, newValues[j]).Scan(&row_id)
		if err != nil {
			fmt.Print("ERRor 1 ")
			fmt.Println(err)
		}
	}
}
func uploadJSONdataToDB() {
	file, _ := ioutil.ReadFile("icecream.json")
	var newProduct []models.Product
	_ = json.Unmarshal([]byte(file), &newProduct)
	insertIntoDB(newProduct)

	/*for i := 0; i < len(newProduct); i++ {
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
	}*/

}

func main() {
	db = driver.ConnectDB()
	//uploadJSONdataToDB()
	router := mux.NewRouter()

	controller := controllers.Controller{}

	docs.SwaggerInfo.Title = "BEN JERRY API v1"

	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")

	router.HandleFunc("/protectedEndpoint", utils.TokenVerifyMiddleWare(controller.Protected(db))).Methods("GET")

	router.HandleFunc("/products", utils.TokenVerifyMiddleWare(controller.GetProducts(db))).Methods("GET")
	router.HandleFunc("/products/{id}", utils.TokenVerifyMiddleWare(controller.GetProduct(db))).Methods("GET")
	//router.HandleFunc("/products/{id}", controller.GetProduct(db)).Methods("GET")
	router.HandleFunc("/products", utils.TokenVerifyMiddleWare(controller.AddProduct(db))).Methods("POST")
	router.HandleFunc("/products/{id}", utils.TokenVerifyMiddleWare(controller.UpdateProduct(db))).Methods("PUT")
	router.HandleFunc("/products/{id}", utils.TokenVerifyMiddleWare(controller.RemoveProduct(db))).Methods("DELETE")

	router.HandleFunc("/products/{id}/name", utils.TokenVerifyMiddleWare(controller.GetProductName(db))).Methods("GET")
	router.HandleFunc("/products/{id}/name", utils.TokenVerifyMiddleWare(controller.UpdateProductName(db))).Methods("PUT")

	router.HandleFunc("/products/{id}/image_open", utils.TokenVerifyMiddleWare(controller.GetProductImageOpen(db))).Methods("GET")
	router.HandleFunc("/products/{id}/image_open", utils.TokenVerifyMiddleWare(controller.UpdateProductImageOpen(db))).Methods("PUT")

	router.HandleFunc("/products/{id}/image_closed", utils.TokenVerifyMiddleWare(controller.GetProductImageClose(db))).Methods("GET")
	router.HandleFunc("/products/{id}/image_closed", utils.TokenVerifyMiddleWare(controller.UpdateProductImageClose(db))).Methods("PUT")

	router.HandleFunc("/products/{id}/description", utils.TokenVerifyMiddleWare(controller.GetProductDescription(db))).Methods("GET")
	router.HandleFunc("/products/{id}/description", utils.TokenVerifyMiddleWare(controller.UpdateProdutDdescription(db))).Methods("PUT")

	router.HandleFunc("/products/{id}/story", utils.TokenVerifyMiddleWare(controller.GetProductStory(db))).Methods("GET")
	router.HandleFunc("/products/{id}/story", utils.TokenVerifyMiddleWare(controller.UpdateProductStory(db))).Methods("PUT")

	router.HandleFunc("/products/{id}/allergy", utils.TokenVerifyMiddleWare(controller.GetProductAllergy(db))).Methods("GET")
	router.HandleFunc("/products/{id}/allergy", utils.TokenVerifyMiddleWare(controller.UpdateProductAllergy(db))).Methods("PUT")

	router.HandleFunc("/products/{id}/diet", utils.TokenVerifyMiddleWare(controller.GetProductDiet(db))).Methods("GET")
	router.HandleFunc("/products/{id}/diet", utils.TokenVerifyMiddleWare(controller.UpdateProductDiet(db))).Methods("PUT")

	router.HandleFunc("/products/{id}/ingredients", utils.TokenVerifyMiddleWare(controller.GetProductIngredients(db))).Methods("GET")
	router.HandleFunc("/products/{id}/ingredients", utils.TokenVerifyMiddleWare(controller.UpdateProductIngredients(db))).Methods("PUT")

	router.HandleFunc("/products/{id}/sourcingvalue", utils.TokenVerifyMiddleWare(controller.GetProductSourcingValues(db))).Methods("GET")
	router.HandleFunc("/products/{id}/sourcingvalue", utils.TokenVerifyMiddleWare(controller.UpdateProductSourcingValues(db))).Methods("PUT")

	//router.HandleFunc("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)).Methods("GET")
	//router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler(swaggerFiles.Handler))
	//router.PathPrefix("/").Handler(httpSwagger.WrapHandler)
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8000", router))
}
