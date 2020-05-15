package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type Controller struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c Controller) GetProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Method invoked : SUCCESS")
	}
}
