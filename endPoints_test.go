package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/*var controller controllers.Controller

var db *sql.DB

func init() {
	gotenv.Load()
	db = driver.ConnectDB()
	controller = controllers.Controller{}
}*/
func TestGetEntryByID(t *testing.T) {
	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("id", "4")
	req.Header.Add("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im5pcm1hbDVAYmVuamVycnkuY29tIiwiaXNzIjoiY291cnNlIn0.KOs1A1nFYjkyCQiOMdZngjc16ML8oIvx__oTf3UQLMI")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	t.Errorf(rr.Body.String())
	expected := `{"id":4 }` //,"first_name":"Krish","last_name":"Bhanushali","email_address":"krishsb2405@gmail.com","phone_number":"0987654321"}`
	if rr.Body.String() != expected {
		//t.Errorf("handler returned unexpected body: got %v want %v",
		//	rr.Body.String(), expected)
	}
}
