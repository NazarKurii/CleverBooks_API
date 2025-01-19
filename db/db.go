package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not open the database")
	}
	createTables()

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

}

func createTables() {
	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			registered BOOLEAN
		)
	`
	_, err := DB.Exec(createUsersTable)

	if err != nil {
		panic(err.Error())
	}

	createFavotitesTable := `
		CREATE TABLE IF NOT EXISTS favorites (
		user_id INTEGER NOT NULL,
		book_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
		FOREIGN KEY (book_id) REFERENCES catalogue(id)
		)
	`
	_, err = DB.Exec(createFavotitesTable)

	if err != nil {
		panic(err.Error())
	}

	createCartsTable := `
		CREATE TABLE IF NOT EXISTS carts (
		user_id INTEGER NOT NULL,
		book_id INTEGER NOT NULL,
		amount INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
		FOREIGN KEY (book_id) REFERENCES catalogue(id)
		)
	`
	_, err = DB.Exec(createCartsTable)

	if err != nil {
		panic(err.Error())
	}

	createBooksTable := `
		CREATE TABLE IF NOT EXISTS catalogue (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		genre TEXT NOT NULL,
		year INTEGER NOT NULL,
		language TEXT NOT NULL,
		price FLOAT NOT NULL,
		booksLeft INTEGER NOT NULL,
		eBook BOOLEAN NOT NULL,
		audioBook BOOLEAN NOT NULL,
		book BOOLEAN NOT NULL,
		imageURL TEXT NOT NULL
		)
	`
	_, err = DB.Exec(createBooksTable)

	if err != nil {
		panic(err.Error())
	}

}
