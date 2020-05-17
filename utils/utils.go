package utils

import (
	"ben-jerry/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// RespondWithError error
func RespondWithError(w http.ResponseWriter, status int, error models.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

// ResponseJSON wiht error
func ResponseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GenerateToken wiht error
func GenerateToken(user models.User) (string, error) {
	var err error
	secret := os.Getenv("SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "course",
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

// TokenVerifyMiddleWare .... dfadf
func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObject models.Error
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}

				return []byte(os.Getenv("SECRET")), nil
			})

			if error != nil {
				errorObject.Message = error.Error()
				RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}

			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				errorObject.Message = error.Error()
				RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Message = "Invalid token."
			RespondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}

func getIngredientsIndex(db *sql.DB, values []string) []int {
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
func getSourceValueIndex(db *sql.DB, values []string) []int {
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
func insertIntoDB(db *sql.DB, products []models.Product) {
	for i := 0; i < len(products); i++ {
		insertStringValuesToDB(db, products[i])
		insertIngredients(db, products[i].ProductId, products[i].Ingredients)
		insertSourceValue(db, products[i].ProductId, products[i].SourcingValues)
	}
}
func insertStringValuesToDB(db *sql.DB, product models.Product) {
	id := 0
	stmt := "insert into products (id, name, image_open, image_closed, description, story, allergy_info, dietary_certifications) values($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;"
	err := db.QueryRow(stmt, product.ProductId, product.Name, product.ImageOpen, product.ImageClosed, product.Description, product.Story, product.AllergyInfo, product.DietaryCertifications).Scan(&id)
	if err != nil {
		fmt.Print("ERRor 3 ")
		fmt.Println(err)
	}
}
func insertSourceValue(db *sql.DB, id string, values []string) {
	stmt := "insert into sourcing_values (product_id, value_id) values($1, $2) RETURNING id;"
	row_id := 0
	values = validateParenthesis(values)
	//fmt.Print("Validate parenthses")
	//fmt.Println(values)
	newValues := getSourceValueIndex(db, values)
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
func insertIngredients(db *sql.DB, id string, values []string) {
	row_id := 0
	stmt := "insert into ingredients (product_id, value_id) values($1, $2) RETURNING id;"
	values = validateParenthesis(values)
	newValues := getIngredientsIndex(db, values)
	//fmt.Println(newValues)
	for j := 0; j < len(newValues); j++ {
		err := db.QueryRow(stmt, id, newValues[j]).Scan(&row_id)
		if err != nil {
			fmt.Print("ERRor 1 ")
			fmt.Println(err)
		}
	}
}
func UploadJSONdataToDB(db *sql.DB) {
	file, _ := ioutil.ReadFile("icecream.json")
	var newProduct []models.Product
	_ = json.Unmarshal([]byte(file), &newProduct)
	insertIntoDB(db, newProduct)
}
