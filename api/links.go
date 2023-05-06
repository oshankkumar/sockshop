package api

import "fmt"

var entitymap = map[string]string{
	"customer": "customers",
	"address":  "addresses",
	"card":     "cards",
}

type Links map[string]Href

type Href struct {
	Href string `json:"href"`
}

func (l *Links) AddLink(domain string, ent string, id string) {
	nl := make(Links)
	link := fmt.Sprintf("http://%v/%v/%v", domain, entitymap[ent], id)
	nl[ent] = Href{link}
	nl["self"] = Href{link}
	*l = nl
}

func (l *Links) AddAttrLink(domain string, attr string, corent string, id string) {
	link := fmt.Sprintf("http://%v/%v/%v/%v", domain, entitymap[corent], id, entitymap[attr])
	nl := *l
	nl[entitymap[attr]] = Href{link}
	*l = nl
}

func (l *Links) AddCustomer(domain string, id string) {
	l.AddLink(domain, "customer", id)
	l.AddAttrLink(domain, "address", "customer", id)
	l.AddAttrLink(domain, "card", "customer", id)
}

func (l *Links) AddAddress(domain string, id string) {
	l.AddLink(domain, "address", id)
}

func (l *Links) AddCard(domain string, id string) {
	l.AddLink(domain, "card", id)
}
