package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"test/db"
)

type Book struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Genre     string  `json:"genre"`
	Year      int     `json:"year"`
	Language  string  `json:"language"`
	Price     float64 `json:"price"`
	BooksLeft int     `json:"booksLeft"`
	EBook     bool    `json:"eBook"`
	AudioBook bool    `json:"audioBook"`
	Book      bool    `json:"book"`
	ImageURL  string  `json:"imageURL"`
	Favorite  bool    `json:"favorite"`
	Cart      int     `json:"cart"`
	IsNew     bool    `json:"isNew"`
}

type Catalogue []Book

func (catalogue *Catalogue) GetBooksInfo(IDs []int, userID int64) error {

	placeHolders, args, err := createPlaceHolders(IDs)

	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
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
	WHERE catalogue.id IN (%s) 
	;
	`, strings.Join(placeHolders, ", "))

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	args = append([]interface{}{userID, userID}, args...)
	rows, err := stmt.Query(args...)

	if err != nil {
		return err
	}

	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Language, &book.Price, &book.BooksLeft, &book.EBook, &book.AudioBook, &book.Book, &book.ImageURL, &book.Favorite, &book.Cart)

		if err != nil {
			return err
		}

		book.IsNew = true

		*catalogue = append(*catalogue, book)
	}

	return err
}

func (c Catalogue) Sort() map[string]Catalogue {
	var catalogues = map[string]Catalogue{}

	for _, book := range c {
		catalogues[book.Genre] = append(catalogues[book.Genre], book)
	}

	return catalogues
}

func CreateTemporaryCatalogue() {

	row := db.DB.QueryRow("SELECT id FROM catalogue LIMIT 1 ")

	var id int = 0
	_ = row.Scan(&id)

	if id == 1 {
		return
	}

	catalogueFile, err := os.Open("./storage.json")

	if err != nil {
		panic(err.Error())
	}

	defer catalogueFile.Close()

	readFile, err := io.ReadAll(catalogueFile)

	if err != nil {
		panic(err.Error())
	}

	var catalogue []Book

	json.Unmarshal(readFile, &catalogue)

	stmt, err := db.DB.Prepare(`
		INSERT INTO catalogue (title, author, genre, year, language, price, booksLeft, eBook, audioBook, book, imageURL)
			Values(?,?,?,?,?,?,?,?,?,?,?)
	`)

	if err != nil {
		panic(err.Error())
	}

	for _, book := range catalogue {
		_, err := stmt.Exec(book.Title, book.Author, book.Genre, book.Year, book.Language, book.Price, book.BooksLeft, book.EBook, book.AudioBook, book.Book, book.ImageURL)
		if err != nil {
			panic(err.Error())
		}
	}

	var user User
	token, _ := user.CreateGuest()

	fmt.Println(token)
}

func createPlaceHolders[T any](elements []T) ([]string, []interface{}, error) {

	length := len(elements)

	if length == 0 {
		return nil, nil, errors.New("No IDs were provided")
	}

	var placeHolders = make([]string, length)
	args := make([]interface{}, length)

	for i := 0; i < length; i++ {
		placeHolders[i] = "?"
		args[i] = elements[i]
	}

	return placeHolders, args, nil
}
