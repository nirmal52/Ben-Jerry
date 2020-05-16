package main

import (
	"ben-jerry/controllers"
	"ben-jerry/driver"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

/*var controller controllers.Controller

var db *sql.DB

func init() {
	gotenv.Load()
	db = driver.ConnectDB()
	controller = controllers.Controller{}
}*/

func setUpTadfest() {
	db = driver.ConnectDB()
	router := mux.NewRouter()
	controller := controllers.Controller{}
	router.HandleFunc("/endpoint", controller.Protected(db)).Methods("GET")
	log.Println("inside setup!")

	log.Fatal(http.ListenAndServe(":8000", router))
	log.Println("listen set")
}
func TestdddMain(m *testing.M) {
	setUpTest()
	log.Println("SDDD")
	exitVal := m.Run()
	log.Println("Do stuff AFTER the tests!")

	os.Exit(exitVal)
}

func TestA(t *testing.T) {
	log.Println("TestA running")
}

func TestB(t *testing.T) {
	log.Println("TestB running")
}

/*
func Router() *mux.Router {
	gotenv.Load()
	var db *sql.DB
	db = driver.ConnectDB()
	//uploadJSONdataToDB()
	router := mux.NewRouter()

	controller := controllers.Controller{}

	router.HandleFunc("/endpoint", controller.Protected(db)).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
	return router
}

func TestCreateEndpoint(t *testing.T) {
	request, _ := http.NewRequest("GET", "/endpoint", nil)
	response := httptest.NewRecorder()
	fmt.Println("HERE")
	Router().ServeHTTP(response, request)
	fmt.Println("HERE 2 ")
	assert.Equal(t, 200, response.Code, "OK response is expected")
}
*/
