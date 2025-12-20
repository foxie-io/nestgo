package models

type Order struct {
	ID       int
	UserID   int
	Product  string
	Quantity int
}
