package models

import (
	"errors"

	"test/db"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type Order struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	PhoneNumber     string `json:"phoneNumber"`
	DeliveryCompany string `json:"deliveryCompany"`
	Comment         string `json:"comment"`
	Items           []Book `json:"items"`
	UserID          int64  `json:"userId"`
	Adress          Adress `json:"adress"`
	SessionID       string `json:"sessionId"`
	Status          string `json:"status"`
}

const (
	StatusWaitingForPayment  = "Waiting for payment"
	StatusPaymentSuccessful  = "Payment successful"
	StatusPaymentFailed      = "Payment failed"
	StatusDelivered          = "Delivered"
	StatusOnItsWayToCustomer = "On its way to you"
)

type Orders []Order

var CouldNotRegisterOrderError = errors.New("Payment successful, but order was not registered")

func (order *Order) Payment() (string, error) {

	stripe.Key = "rk_test_51QoptdJZ010pKpoPII2ZLss9dvJtMfwy8o0L6ev4PWeJJZQhWq6jTeIR4jhqnAcCVfSSUdMEeBl95lBa8HeosQso009EtmzbUn"

	lineItems, err := cartItemsToStripeitems(order.Items)

	if err != nil {
		return "", err
	}

	domain := "http://localhost:3000/checkout"
	params := &stripe.CheckoutSessionParams{
		LineItems:  lineItems,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		CancelURL:  stripe.String(domain + "/failed?session_id={CHECKOUT_SESSION_ID}"),
		SuccessURL: stripe.String(domain + "/successful?session_id={CHECKOUT_SESSION_ID}"),
	}

	s, err := session.New(params)

	if err != nil {
		return "", err
	}

	order.SessionID = s.ID

	return s.URL, err
}

func cartItemsToStripeitems(books []Book) ([]*stripe.CheckoutSessionLineItemParams, error) {
	stripeItems := make([]*stripe.CheckoutSessionLineItemParams, len(books))

	for i, book := range books {
		id, err := book.GetStripeID()
		if err != nil {
			return nil, err
		}
		stripeItems[i] = &stripe.CheckoutSessionLineItemParams{

			Price:    stripe.String(id),
			Quantity: stripe.Int64(int64(book.Cart)),
		}
	}

	return stripeItems, nil
}

func (order Order) CreateOrder() error {
	query := "INSERT into orders (name, phone_number, delivery_company, comment, user_id, country, city, street, house_number, flat_number, post_code, session_id, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(order.Name, order.PhoneNumber, order.DeliveryCompany, order.Comment, order.UserID, order.Adress.Country, order.Adress.City, order.Adress.Street, order.Adress.HouseNumber, order.Adress.FlatNumber, order.Adress.PostCode, order.SessionID, StatusWaitingForPayment)

	if err != nil {
		return err
	}

	if effected, _ := result.RowsAffected(); effected == 0 {
		return db.NoRowsEffectedError
	}

	order.ID, err = result.LastInsertId()

	if err != nil {
		return err
	}

	err = order.saveOrderItems()

	return err
}

func (order Order) saveOrderItems() error {

	query := "INSERT INTO ordered_items (book_id, order_id) VALUES (?,?)"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	for _, item := range order.Items {
		result, err := stmt.Exec(item.ID, order.ID)
		if err != nil {
			return err
		}

		rowsEffected, _ := result.RowsAffected()

		if rowsEffected == 0 {
			return db.NoRowsEffectedError
		}
	}

	return nil

}

func (order Order) PaymentSuccessful() error {

	query := "UPDATE orders SET status = ? WHERE session_id = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(StatusPaymentSuccessful, order.SessionID)

	if effected, _ := result.RowsAffected(); effected == 0 {
		return db.NoRowsEffectedError
	}

	return err

}

func (order Order) PaymentFailed() error {

	query := "UPDADE orders SET status = ? WHERE session_id = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(StatusPaymentFailed, order.SessionID)

	if effected, _ := result.RowsAffected(); effected == 0 {
		return db.NoRowsEffectedError
	}

	return err

}

func GetOrders(userID int64) ([]Order, error) {
	query := "SELECT * FROM orders WHERE user_id = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userID)

	if err != nil {
		return nil, err
	}

	var orders Orders

	for rows.Next() {
		var order Order

		err = rows.Scan(&order.ID, &order.SessionID, &order.Name, &order.PhoneNumber, &order.DeliveryCompany, &order.Comment, &order.UserID, &order.Adress.Street, &order.Adress.City, &order.Adress.PostCode, &order.Adress.FlatNumber, &order.Adress.HouseNumber, &order.Adress.Country, &order.Status)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	err = orders.getItems()

	return orders, err
}

func (orders *Orders) getItems() error {

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
    	catalogue.imageURL
	FROM catalogue
	LEFT JOIN ordered_items 
		on ordered_items.book_id = catalogue.id 
	WHERE ordered_items.order_id = ?
	;
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	ordersSlice := []Order(*orders)

	for i, order := range ordersSlice {

		rows, err := stmt.Query(order.ID)

		if err != nil {
			return err
		}

		for rows.Next() {
			var book Book
			err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Language, &book.Price, &book.BooksLeft, &book.EBook, &book.AudioBook, &book.Book, &book.ImageURL)

			if err != nil {
				return err
			}

			ordersSlice[i].Items = append(ordersSlice[i].Items, book)

		}

	}

	*orders = ordersSlice

	return nil
}
