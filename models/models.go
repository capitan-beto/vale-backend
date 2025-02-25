package models

import "github.com/mercadopago/sdk-go/pkg/preference"

type ContestantData struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Number    int    `json:"number"`
	CreatedAt string
	Confirmed bool   `json:"confirmed"`
	ExtRef    string `json:"ext_ref"`
}

type WebhookPayload struct {
	Id map[string]string `json:"data"`
}

type PaymentData struct {
	ExternalReference string                 `json:"ext_ref"`
	Item              preference.ItemRequest `json:"item"`
}
