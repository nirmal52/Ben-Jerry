package userRepository

import (
	"ben-jerry/models"
	"database/sql"
	"fmt"
)

type UserRepository struct{}

func logError(err error, message string) {
	if err != nil {
		fmt.Print(message)
		fmt.Println(err)
	}
}

func (u UserRepository) Signup(db *sql.DB, user models.User) (models.User, error) {

	stmt := "insert into users (email, password) values($1, $2) RETURNING id;"
	err := db.QueryRow(stmt, user.Email, user.Password).Scan(&user.ID)

	if err != nil {
		logError(err, "SignUp")
		return user, err
	}

	user.Password = ""
	return user, nil
}

func (u UserRepository) Login(db *sql.DB, user models.User) (models.User, error) {
	row := db.QueryRow("select * from users where email=$1", user.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		logError(err, "Login")
		return user, err
	}

	return user, nil
}
