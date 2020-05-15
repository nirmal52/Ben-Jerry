package main

import (
	"database/sql"
	"fmt"
	"jwt-course/driver"
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
func main() {
	db = driver.ConnectDB()
	router := mux.NewRouter()

	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
