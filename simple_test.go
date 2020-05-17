package main

import (
	"ben-jerry/controllers"
	"ben-jerry/driver"
	"ben-jerry/models"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var router *mux.Router
var ts *httptest.Server

func setUpTest() {
	db = driver.ConnectDB()
	router = mux.NewRouter()
	controller := controllers.Controller{}

	router.HandleFunc("/products", controller.GetProducts(db)).Methods("GET")
	router.HandleFunc("/products/{id}", controller.GetProduct(db)).Methods("GET")
	router.HandleFunc("/products", controller.AddProduct(db)).Methods("POST")
	router.HandleFunc("/products/{id}", controller.UpdateProduct(db)).Methods("PUT")
	router.HandleFunc("/products/{id}", controller.RemoveProduct(db)).Methods("DELETE")
	ts = httptest.NewServer(router)
}

func TestMain(m *testing.M) {
	setUpTest()
	defer ts.Close()
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
/*func bookmarkIndex(t *testing.T) {
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

func productName(t *testing.T) {
	//log.Println("Name")
	ts := httptest.NewServer(router)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/products/2189/name")
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	newStr := buf.String()
	//log.Println(newStr)
	expected := "Chillin' the Roast"
	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}
	assert.Equal(t, expected, newStr, "OK response is expected")
}
*/
//var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
//req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

func createNewProduct() models.Product {
	var newProduct models.Product
	newProduct.Name = "Creating new Product to test Post"
	newProduct.ImageClosed = "NEW image open"
	newProduct.ImageOpen = "NEW image open"
	newProduct.Description = "Description"
	newProduct.AllergyInfo = "Allergy"
	newProduct.DietaryCertifications = "Diet"
	newProduct.Ingredients = []string{"ing1", "ing3"}
	newProduct.SourcingValues = []string{"sv1", "sv2"}
	newProduct.Story = "Story "
	newProduct.ProductId = "626"
	return newProduct
}

/*
	Test Run:
	1. HTTP POST - Create a new product
	2. Save response of productID
	3. HTTP GET -  Try to get the product by ID
	4. Compare both the object of POST and GET
*/
func Test_AddandCompareProduct(t *testing.T) {
	product := createNewProduct()
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(product)
	res, err := http.Post(ts.URL+"/products", "application/json", buf)

	if err != nil {
		log.Println(err)
	}

	resbuf := new(bytes.Buffer)
	resbuf.ReadFrom(res.Body)
	response := resbuf.String()

	jsonMap := make(map[string]int)
	err = json.Unmarshal([]byte(response), &jsonMap)
	if err != nil {
		log.Println(err)
	}

	createdID := strconv.Itoa(jsonMap["productID"])
	product.ProductId = createdID
	expected := product

	getResponse, err := http.Get(ts.URL + "/products/" + product.ProductId)
	if err != nil {
		log.Println(err)
	}

	getresbuf := new(bytes.Buffer)
	getresbuf.ReadFrom(getResponse.Body)
	productCreated := getresbuf.String()

	var getProduct models.Product
	json.Unmarshal([]byte(productCreated), &getProduct)

	if getResponse.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}
	assert.Equal(t, expected, getProduct, "OK")
}

/*
func Test_PutProductName(t *testing.T) {
	//log.Println("Name")
	ts := httptest.NewServer(router)
	defer ts.Close()
	client := &http.Client{}
	product := createNewProduct()
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(product)
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/products/626", buf)
	res, err := client.Do(req)
	resbuf := new(bytes.Buffer)
	resbuf.ReadFrom(res.Body)
	response := resbuf.String()
	log.Println(response)
	expected := "Chillin' the Roast"
	if err != nil {
		t.Errorf("Expected nil, received %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}
	assert.Equal(t, expected, response, "OK response is expected")
}
*/
