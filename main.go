package main

import (
	"ben-jerry/controllers"
	"ben-jerry/docs"
	"ben-jerry/driver"
	"ben-jerry/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	//"github.com/swaggo/gin-swagger/swaggerFiles"
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
// @BasePath ""

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

	router.HandleFunc("/products/{id}/name", utils.TokenVerifyMiddleWare(controller.GetProductDetail(db, 0))).Methods("GET")
	router.HandleFunc("/products/{id}/name", utils.TokenVerifyMiddleWare(controller.UpdateProductDetail(db, 0))).Methods("PUT")

	router.HandleFunc("/products/{id}/image_open", utils.TokenVerifyMiddleWare(controller.GetProductDetail(db, 1))).Methods("GET")
	router.HandleFunc("/products/{id}/image_open", utils.TokenVerifyMiddleWare(controller.UpdateProductDetail(db, 1))).Methods("PUT")

	router.HandleFunc("/products/{id}/image_closed", utils.TokenVerifyMiddleWare(controller.GetProductDetail(db, 2))).Methods("GET")
	router.HandleFunc("/products/{id}/image_closed", utils.TokenVerifyMiddleWare(controller.UpdateProductDetail(db, 2))).Methods("PUT")

	router.HandleFunc("/products/{id}/description", utils.TokenVerifyMiddleWare(controller.GetProductDetail(db, 3))).Methods("GET")
	router.HandleFunc("/products/{id}/description", utils.TokenVerifyMiddleWare(controller.UpdateProductDetail(db, 3))).Methods("PUT")

	router.HandleFunc("/products/{id}/story", utils.TokenVerifyMiddleWare(controller.GetProductDetail(db, 4))).Methods("GET")
	router.HandleFunc("/products/{id}/story", utils.TokenVerifyMiddleWare(controller.UpdateProductDetail(db, 4))).Methods("PUT")

	router.HandleFunc("/products/{id}/allergy", utils.TokenVerifyMiddleWare(controller.GetProductDetail(db, 5))).Methods("GET")
	router.HandleFunc("/products/{id}/allergy", utils.TokenVerifyMiddleWare(controller.UpdateProductDetail(db, 5))).Methods("PUT")

	router.HandleFunc("/products/{id}/diet", utils.TokenVerifyMiddleWare(controller.GetProductDetail(db, 6))).Methods("GET")
	router.HandleFunc("/products/{id}/diet", utils.TokenVerifyMiddleWare(controller.UpdateProductDetail(db, 6))).Methods("PUT")

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

/*
func main() {
	db = driver.ConnectDB()
	//uploadJSONdataToDB()
	r := gin.Default()

	c := controllers.Controller{}

	docs.SwaggerInfo.Title = "BEN JERRY API v1"

	v1 := r.Group("api/v1")
	{
		products := v1.Group("/products")
		{
			//products.GET(":id", utils.TokenVerifyMiddleWare(conProtected(db)) )
			products.GET("", utils.TokenVerifyMiddleWare(c.GetProducts(db)))
			/*products.POST("", c.AddAccount)
			products.DELETE(":id", c.DeleteAccount)
			products.PATCH(":id", c.UpdateAccount)
			products.POST(":id/images", c.UploadAccountImage)*/
/*	}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
///}
*/
