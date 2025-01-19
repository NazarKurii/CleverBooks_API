package models

import (
	"errors"
	"fmt"
	"test/db"
	"test/utils"
)

type User struct {
	ID         int64  `json:"id"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Registered bool   `json:"registered"`
}

func (u User) Save(token string) error {
	var query string

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	var args []interface{}
	if token == "" {
		query = "INSERT INTO users (email, password) VALUES (?,?)"

		args = append(args, u.Email, hashedPassword)
	} else {
		id, err := utils.VerifyToken(token)

		if err != nil {
			return err
		}

		query = fmt.Sprintf("UPDATE users SET email = ?, password = ? WHERE id = ?")

		args = append(args, u.Email, hashedPassword, id)
	}

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(args...)

	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()
	u.ID = userID

	return err
}

func (u *User) ValidateCredentials() error {
	query := `SELECT password, id FROM users WHERE email = ?`

	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&retrievedPassword, &u.ID)

	if err != nil {
		return err
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)

	if !passwordIsValid {
		return errors.New("Credentials invalid")
	}

	return nil
}

func (u *User) CreateGuest() error {
	query := "SELECT id FROM users ORDER BY id DESC LIMIT 1"

	row := db.DB.QueryRow(query)

	var ID int64

	err := row.Scan(ID)

	if err != nil {
		return err
	}

	ID++

	u.Email = fmt.Sprintf("email_%v", ID)
	u.Password = fmt.Sprintf("password%v", ID)

	return nil
}
