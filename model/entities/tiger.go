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

type Sighting struct {
	ID          uint      `json:"id"`
	TigerID     uint      `json:"tigerId" gorm:"foreignKey:TigerID;references:ID"`
	UserID      uint      `json:"userId" gorm:"foreignKey:UserID;references:ID"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	SeenAt      time.Time `json:"seenAt"`
	ImageURL    string    `json:"imageURL"`
	Description string    `json:"description"`
	UpdatedAt   time.Time
	CreatedAt   time.Time
	Tiger       Tiger
	User        User
}
