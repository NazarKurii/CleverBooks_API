package models

import (
	"fmt"
	"test/db"
)

type Favorite struct {
	UserID     int64 `json:"userId" binding:"required"`
	FavoriteID int64 `json:"favoriteId" binding:"required"`
}

func (f Favorite) Save() error {
	query := "INSERT INTO favorites (user_id, book_id) VALUES (?, ?)"
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
	query := "DELETE FROM favorites WHERE user_id = ? AND book_id = ?"
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

func GetFavorites(userID int64) ([]Book, error) {
	query := `SELECT 
    	catalogue.id,
    	catalogue.title,
    	catalogue.author,
    	catalogue.genre,
    	catalogue.year,
    	catalogue.language,
    	catalogue.price,
    	catalogue.booksLeft,
    	catalogue.eBook,
    	catalogue.audioBook,
    	catalogue.book,
    	catalogue.imageURL,
    	CASE WHEN favorites.book_id IS NOT NULL THEN 1 ELSE 0 END AS favorite,
   		COALESCE(carts.amount, 0) AS cart
	FROM catalogue
	LEFT JOIN carts
		on carts.book_id = catalogue.id AND  carts.user_id = ?  
	LEFT JOIN favorites 
		on favorites.book_id = catalogue.id AND favorites.user_id = ?
	WHERE favorites.user_id == ?
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return nil, fmt.Errorf("failed to prepare getFavorites query: %w", err)
	}

	defer stmt.Close()

	result, err := stmt.Query(userID, userID, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to execute getFavorites query: %w", err)
	}

	defer result.Close()

	var favorites []Book
	for result.Next() {
		var book Book
		err := result.Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Language, &book.Price, &book.BooksLeft, &book.EBook, &book.AudioBook, &book.Book, &book.ImageURL, &book.Favorite, &book.Cart)

		if err != nil {
			return nil, fmt.Errorf("failed to scan favorite ID: %w", err)
		}

		favorites = append(favorites, book)
	}

	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error during result iteration: %w", err)
	}

	return favorites, nil
}
