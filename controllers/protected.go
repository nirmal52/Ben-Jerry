package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (c Controller) Protected(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("SUCCESS")
	}
}
