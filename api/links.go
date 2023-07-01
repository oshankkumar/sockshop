package api

import "fmt"

type Links map[string]Href

type Href struct {
	Href string `json:"href"`
}

func NewCustomerLinks(domain string, id string) Links {
	l := make(Links)
	l["self"] = Href{fmt.Sprintf("http://%v/customers/%v", domain, id)}
	l["customer"] = Href{fmt.Sprintf("http://%v/customers/%v", domain, id)}
	l["addresses"] = Href{fmt.Sprintf("http://%v/customers/%v/addresses", domain, id)}
	l["cards"] = Href{fmt.Sprintf("http://%v/customers/%v/cards", domain, id)}
	return l
}

func NewAddressLinks(domain string, id string) Links {
	l := make(Links)
	l["self"] = Href{fmt.Sprintf("http://%v/addresses/%v", domain, id)}
	l["address"] = Href{fmt.Sprintf("http://%v/addresses/%v", domain, id)}
	return l
}

func NewCardLinks(domain string, id string) Links {
	l := make(Links)
	l["self"] = Href{fmt.Sprintf("http://%v/cards/%v", domain, id)}
	l["card"] = Href{fmt.Sprintf("http://%v/cards/%v", domain, id)}
	return l
}
