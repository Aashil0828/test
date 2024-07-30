package models

// PersonRequest represents the request body for creating a new person
type PersonRequest struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Age         int64  `json:"age"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	State       string `json:"state"`
	Street1     string `json:"street1"`
	Street2     string `json:"street2"`
	ZipCode     string `json:"zip_code"`
}
