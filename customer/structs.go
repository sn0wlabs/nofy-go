package customer

import (
	"time"
)

type Customer struct {
	ID          string    `json:"id"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	RelatedUser string    `json:"relatedUser"`
	CustomerID  string    `json:"customerId"`
	Billing     *Address  `json:"billing"`
	Shipping    *Address  `json:"shipping"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
}

type Address struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Company   string `json:"company"`
	Street    string `json:"street"`
	Street2   string `json:"street2"`
	HouseNr   string `json:"houseNr"`
	Zip       string `json:"zip"`
	City      string `json:"city"`
	Country   string `json:"country"`
}
