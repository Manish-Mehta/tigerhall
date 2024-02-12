package dto

import "time"

type Coordinate struct {
	Lat float64 `json:"lat" binding:"required"`
	Lon float64 `json:"lon" binding:"required"`
}

type TigerCreateRequest struct {
	Name       string     `json:"name" binding:"required"`
	DOB        string     `json:"dob" binding:"required"`
	LastSeen   time.Time  `json:"lastSeen" binding:"required"`
	Coordinate Coordinate `json:"coordinate" binding:"required"`
}

type TigerListResponse struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}
