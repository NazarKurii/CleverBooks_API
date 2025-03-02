package models

import (
	"fmt"
	"test/db"
)

type Adress struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	Street      string `json:"street"`
	City        string `json:"city"`
	PostCode    string `json:"postCode"`
	FlatNumber  string `json:"flatNumber"`
	HouseNumber string `json:"houseNumber"`
	Country     string `json:"country"`
}

func GetAdresses(userID int64) ([]Adress, error) {
	query := "SELECT id,user_id, street, city, post_code, flat_number, house_number, country FROM adresses WHERE user_id = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var adresses = []Adress{}

	for rows.Next() {
		var adress Adress

		err := rows.Scan(&adress.ID, &adress.UserID, &adress.Street, &adress.City, &adress.PostCode, &adress.FlatNumber, &adress.HouseNumber, &adress.Country)

		if err != nil {
			return nil, err
		}

		adresses = append(adresses, adress)
	}

	return adresses, nil
}

func (a Adress) Save() (int64, error) {

	var query string
	var args []interface{}
	if a.ID == 0 {
		query = "INSERT INTO adresses (user_id, street, city, post_code, flat_number, house_number, country) VALUES (?, ?, ?, ?, ?, ?, ?)"
		args = append(args, a.UserID, a.Street, a.City, a.PostCode, a.FlatNumber, a.HouseNumber, a.Country)
	} else {
		query = "UPDATE adresses SET street = ?, CITY = ?, POST_CODE = ?, flat_number = ?, house_number = ?, country = ? WHERE id = ?"
		args = append(args, a.Street, a.City, a.PostCode, a.FlatNumber, a.HouseNumber, a.Country, a.ID)
	}
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare save query: %w", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(args...)

	if err != nil {
		return 0, fmt.Errorf("failed to asign an ID: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return 0, fmt.Errorf("no rows were inserted")
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("failed to asign an ID: %w", err)
	}

	return id, nil
}

func (a Adress) Delete() error {
	query := "DELETE FROM adresses WHERE user_id = ? AND id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare delete query: %w", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(a.UserID, a.ID)

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted")
	}

	return err
}
