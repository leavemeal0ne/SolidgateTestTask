package domen

import (
	"fmt"
	"regexp"
)

type Card struct {
	PAN             int `json:"Card number"`
	ExpirationMonth int `json:"Expiration month"`
	ExpirationYear  int `json:"Expiration year"`
}

type CardIssuer struct {
	Issuer     string         `json:"Issuer"`
	CardLength []string       `json:"CardLength"`
	RgxStr     string         `json:"Rgx"`
	Rgx        *regexp.Regexp `json:"-"`
}

type Response struct {
	Valid bool       `json:"valid"`
	Error *ErrorInfo `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func EmptyFieldRaiseErr(card *Card) error {
	if card.PAN == 0 {
		return fmt.Errorf("Card number is required")
	}
	if card.ExpirationMonth == 0 {
		return fmt.Errorf("Expiration month is required")
	}
	if card.ExpirationYear == 0 {
		return fmt.Errorf("ExpirationYear is required")
	}
	return nil
}
