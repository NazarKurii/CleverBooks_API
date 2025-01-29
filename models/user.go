package models

import (
	"database/sql"
	"errors"
	"fmt"
	"test/db"
	"test/utils"
)

type User struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Number     string `json:"number"`
	Password   string `json:"password"`
	Registered bool   `json:"registered"`
}

func (u User) Save() (string, error) {
	query := "UPDATE users SET name = ?, surname = ?, number = ?, email = ?, password = ?, registered = TRUE WHERE id = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return "", err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return "", err
	}

	_, err = stmt.Exec(u.Name, u.Surname, u.Number, u.Email, hashedPassword, u.ID)

	if err != nil {
		return "", err
	}

	token, err := utils.GenerateUserToken(u.Email, u.ID)

	fmt.Println(u)

	return token, err
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

func (u *User) CreateGuest() (string, error) {
	query := "INSERT INTO users (registered) VALUES (FALSE)"

	response, err := db.DB.Exec(query)

	ID, err := response.LastInsertId()

	if err != nil {
		return "", err
	}

	hashedPassword, err := utils.HashPassword(fmt.Sprintf("password%v", ID))

	if err != nil {
		return "", err
	}

	u.Password = hashedPassword
	u.Email = fmt.Sprintf("num_%v", ID)
	u.ID = ID

	query = "UPDATE users SET email = ?, password = ? WHERE id = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return "", err
	}

	defer stmt.Close()

	_, err = stmt.Exec(u.Email, u.Password, u.ID)

	if err != nil {
		return "", err
	}

	token, err := utils.GenerateGuestToken(u.Email, u.ID)

	return token, err
}

func (u User) VerifyEmail() (bool, error) {

	query := "SELECT email FROM users WHERE email = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(u.Email)

	err = row.Scan(&u.Email)

	fmt.Println(u)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func (u *User) IsRegistered() error {

	query := "SELECT registered FROM users WHERE email = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRow(u.Email)

	err = row.Scan(&u.Registered)

	return err
}

func (u *User) GetUser() error {
	query := "SELECT * FROM users WHERE id = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRow(u.ID)

	if err != nil {
		return err
	}

	err = row.Scan(&u.ID, &u.Email, &u.Number, &u.Password, &u.Name, &u.Surname, &u.Registered)

	return err
}
