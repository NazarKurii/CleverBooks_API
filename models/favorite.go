package models

import (
	"fmt"
	"test/db"
)

type Favorite struct {
	UserID     int64 `json:"user_id" binding:"required"`
	FavoriteID int64 `json:"favorite_id" binding:"required"`
}

func (f Favorite) Save() error {
	query := "INSERT INTO favorites (user_id, favorite_id) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare save query: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(f.UserID, f.FavoriteID)
	if err != nil {
		return fmt.Errorf("failed to execute save query: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were inserted")
	}

	return nil
}

func (f Favorite) Delete() error {
	query := "DELETE FROM favorites WHERE user_id = ? AND favorite_id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare delete query: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(f.UserID, f.FavoriteID)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted")
	}

	return nil
}

func GetFavorites(userID int64) ([]int64, error) {
	query := "SELECT favorite_id FROM favorites WHERE user_id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare getFavorites query: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Query(userID)

	if err != nil {
		return nil, fmt.Errorf("failed to execute getFavorites query: %w", err)
	}
	defer result.Close()

	var favorites []int64
	for result.Next() {
		var favoriteID int64
		err := result.Scan(&favoriteID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan favorite ID: %w", err)
		}
		favorites = append(favorites, favoriteID)
	}

	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error during result iteration: %w", err)
	}

	return favorites, nil
}
