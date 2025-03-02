package db

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

var NoRowsEffectedError = errors.New("NO rows were effected")

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
			email TEXT,
			number TEXT,
			password TEXT,
			name TEXT,
			surname TEXT,
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
		imageURL TEXT NOT NULL,
		stripe_id TEXT NOT NULL
		)
	`
	_, err = DB.Exec(createBooksTable)

	if err != nil {
		panic(err.Error())
	}

	createAdresses := `
		CREATE TABLE IF NOT EXISTS adresses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			street TEXT,
  			city TEXT,
			post_code TEXT,
  			flat_number TEXT,
  			house_number TEXT,
  			country TEXT,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`
	_, err = DB.Exec(createAdresses)

	if err != nil {
		panic(err.Error())
	}

	createOrders := `
		CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id TEXT NOT NULL,
			name TEXT NOT NULL,
			phone_number TEXT NOT NULL,
		  	delivery_company TEXT NOT NULL,
			comment TEXT NOT NULL,
  			user_id NUMBER NOT NULL,
  			
			street TEXT NOT NULL,
  			city TEXT NOT NULL,
			post_code TEXT NOT NULL,
  			flat_number TEXT NOT NULL,
  			house_number TEXT NOT NULL,
  			country TEXT NOT NULL,
			status TEXT NOT NULL,
  			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`

	_, err = DB.Exec(createOrders)

	if err != nil {
		panic(err.Error())
	}

	if err != nil {
		panic(err.Error())
	}

	createOrderedItems := `
		CREATE TABLE IF NOT EXISTS ordered_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER NOT NULL,
			book_id INTEGER NOT NULL,
			FOREIGN KEY (order_id) REFERENCES orders(id)
			FOREIGN KEY (book_id) REFERENCES catalogue (id)
	)
`
	_, err = DB.Exec(createOrderedItems)

	if err != nil {
		panic(err.Error())
	}

}
