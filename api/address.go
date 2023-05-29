package api

import "github.com/google/uuid"

type UserAdressesResponse struct {
	Addresses []Address `json:"addresses"`
}

type Address struct {
	ID       uuid.UUID `json:"id"`
	Street   string    `json:"street"`
	Number   string    `json:"number"`
	Country  string    `json:"country"`
	City     string    `json:"city"`
	PostCode string    `json:"postcode"`
	Links    Links     `json:"_links"`
}
