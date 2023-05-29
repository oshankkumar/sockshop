package api

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UserCardsResponse struct {
	Cards []Card `json:"cards"`
}

type Card struct {
	ID      uuid.UUID `json:"id"`
	LongNum string    `json:"longNum"`
	Expires string    `json:"expires"`
	CCV     string    `json:"ccv"`
	Links   Links     `json:"_links"`
}

func (c *Card) MaskCC() {
	l := len(c.LongNum) - 4
	c.LongNum = fmt.Sprintf("%v%v", strings.Repeat("*", l), c.LongNum[l:])
}
