package entities

import "time"

type User struct {
	ID        uint   `json:"id"`
	UserName  string `json:"userName" gorm:"not null"`
	Email     string `json:"email" gorm:"uniqueIndex"`
	Password  string `json:"password" gorm:"not null"`
	UpdatedAt time.Time
	CreatedAt time.Time
}
