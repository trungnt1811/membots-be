package model

import "time"

type AffHistory struct {
	ID        uint      `json:"id"`
	UserId    int       `json:"name"`
	Status    string    `json:"status"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
