package entity

import "time"

type Balance struct {
	ID            int64     `json:"id"`
	Currency      string    `json:"currency"`
	Balance       float64   `json:"balance"`
	FrozenBalance float64   `json:"frozen_balance"`
	UpdatedAt     time.Time `json:"updated_at"`
}
