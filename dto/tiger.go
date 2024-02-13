package dto

import (
	"mime/multipart"
	"time"
)

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
type TigerCreateSightingRequest struct {
	TigerID uint                  `form:"tigerId" binding:"required"`
	UserID  uint                  `form:"userId"`
	Lat     float64               `form:"lat" binding:"required"`
	Lon     float64               `form:"lon" binding:"required"`
	SeenAt  time.Time             `form:"seenAt" binding:"required"`
	Image   *multipart.FileHeader `form:"image" binding:"required"`
}
