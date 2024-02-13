package entities

import (
	"time"

	"gorm.io/datatypes"
)

type Tiger struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name" gorm:"not null;unique"`
	Dob       datatypes.Date `json:"dob" gorm:"not null"`
	LastSeen  time.Time      `json:"lastSeen" gorm:"index:,sort:desc"`
	Lat       float64        `json:"lat"`
	Lon       float64        `json:"lon"`
	UpdatedAt time.Time
	CreatedAt time.Time
}
