package models

import (
	"database/sql"
	"test/db"
)

type CartItem struct {
	UserID int64 `json:"userId"`
	BookID int64 `json:"bookId"`
	Amount int   `json:"amount"`
}

type UserCart struct {
	UserID int64  `json:"userId"`
	Cart   []Book `json:"cart"`
}

func (cart CartItem) AddToCart() error {
	err := cart.GetCartItemAmount()
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	var stmt *sql.Stmt
	if cart.Amount == 0 {
		query := "INSERT INTO carts (amount, user_id, book_id) VALUES (?, ?, ?)"
		stmt, err = db.DB.Prepare(query)

		if err != nil {
			return err
		}

	} else {
		query := "UPDATE carts SET amount = ? WHERE user_id = ? AND book_id = ?"
		stmt, err = db.DB.Prepare(query)

		if err != nil {
			return err
		}

	}

	defer stmt.Close()
	_, err = stmt.Exec(cart.Amount+1, cart.UserID, cart.BookID)

	return err
}

func (cart CartItem) RemoveFromCart() error {
	err := cart.GetCartItemAmount()
	if err != nil {
		return err
	}

	if cart.Amount == 1 {
		return cart.DeleteFromCart()
	}

	query := "UPDATE carts SET amount = ? WHERE user_id = ? AND book_id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(cart.Amount-1, cart.UserID, cart.BookID)

	return err

}

func (cart CartItem) DeleteFromCart() error {

	query := "DELETE FROM carts WHERE user_id = ? AND book_id = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(cart.UserID, cart.BookID)

	return err

}

func (cart *CartItem) GetCartItemAmount() error {

	query := "SELECT amount FROM carts WHERE user_id = ? AND book_id = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRow(cart.UserID, cart.BookID)

	err = row.Scan(&cart.Amount)

	return err

}

func (userCart *UserCart) GetUsersCart() error {
	query := `
	SELECT 
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
	WHERE carts.user_id == ?
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	rows, err := stmt.Query([]interface{}{userCart.UserID, userCart.UserID, userCart.UserID}...)

	if err != nil {
		return err
	}

	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Language, &book.Price, &book.BooksLeft, &book.EBook, &book.AudioBook, &book.Book, &book.ImageURL, &book.Favorite, &book.Cart)

		if err != nil {
			return err
		}

		userCart.Cart = append(userCart.Cart, book)
	}

	return nil
}
