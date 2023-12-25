package entity

import "time"

type Transaction struct {
	ID           int64      `json:"id"`
	Currency     string     `json:"currency"`
	Amount       float64    `json:"amount"`
	WalletOrCard string     `json:"wallet_or_card"`
	Status       string     `json:"status"`
	CreatedAt    *time.Time `json:"created_at"`
}
