package entities

import "time"

type Tiger struct {
	ID        uint   `json:"id"`
	Name      string `json:"Name"`
	UpdatedAt time.Time
	CreatedAt time.Time
}
