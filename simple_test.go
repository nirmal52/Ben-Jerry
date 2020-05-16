package main

import (
	"ben-jerry/controllers"
	"ben-jerry/driver"
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var router *mux.Router

func setUpTest() {
	db = driver.ConnectDB()
	router = mux.NewRouter()
	controller := controllers.Controller{}

	router.Handle("/v1/bookmark", controller.SimpleReturn())
	router.HandleFunc("/endpoint", controller.Protected(db)).Methods("GET")
	router.HandleFunc("/products/{id}/name", controller.GetProductName(db)).Methods("GET")

	log.Println("inside setup!")

	//	log.Fatal(http.ListenAndServe(":8000", router))
	log.Println("listen set")
}

func TestMain(m *testing.M) {
	setUpTest()
	m.Run()
}

/*func Test_bookmarkIndex(t *testing.T) {
	controller := controllers.Controller{}
	r := mux.NewRouter()
	r.Handle("/v1/bookmark", controller.SimpleReturn())
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/v1/bookmark")
	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}
}*/
func Test_bookmarkIndex(t *testing.T) {
	//	controller := controllers.Controller{}
	log.Println("Book mark index")
	ts := httptest.NewServer(router)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/endpoint")
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	newStr := buf.String()
	log.Println(newStr)
	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}
}

func Test_ProductName(t *testing.T) {
	log.Println("Name")
	ts := httptest.NewServer(router)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/products/2189/name")
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	newStr := buf.String()
	log.Println(newStr)
	expected := "Chillin' the Roast"
	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}
	assert.Equal(t, expected, newStr, "OK response is expected")
}
