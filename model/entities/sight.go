package entities

import "time"

type Sight struct {
	ID        uint      `json:"id"`
	TigerID   uint      `json:"tigerId" gorm:"foreignKey:TigerID;references:ID"`
	UserID    uint      `json:"userId" gorm:"foreignKey:UserID;references:ID"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	SeenAt    time.Time `json:"seenAt" gorm:"index:,sort:desc"`
	ImageURL  string    `json:"imageURL"`
	Tiger     Tiger
	User      User
	UpdatedAt time.Time
	CreatedAt time.Time
}
