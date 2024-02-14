package util

import (
	"errors"
	"mime/multipart"
	"time"

	"github.com/Manish-Mehta/tigerhall/api/dto"
	"github.com/Manish-Mehta/tigerhall/internal/config"
)

// Basic Validator for this project
type Validator interface {
	ValDateFormat(string) (time.Time, error)
	valLat(float64) error
	valLon(float64) error
	ValCoord(dto.Coordinate) error
	ValImage(*multipart.FileHeader) (string, error)
}

var NewValidator = func() Validator {
	return &validator{dateFormat: "2006-01-02"}
}

type validator struct {
	dateFormat string
}

func (v *validator) valLat(lat float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	return nil
}

func (v *validator) valLon(lon float64) error {
	if lon < -180 || lon > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	return nil
}

func (v *validator) valImageMime(imgType string) string {
	switch imgType {
	case "image/png":
		return "png"
	case "image/jpeg", "image/jpg":
		return "jpg"
	default:
		return ""
	}
}

func (v *validator) ValCoord(coord dto.Coordinate) error {
	err := v.valLat(coord.Lat)
	if err != nil {
		return err
	}
	return v.valLon(coord.Lon)
}

func (v *validator) ValDateFormat(date string) (time.Time, error) {
	return time.Parse(v.dateFormat, date)
}

func (v *validator) ValImage(img *multipart.FileHeader) (string, error) {
	if img.Size == 0 || img.Size >= config.MAX_UPLOAD_IMAGE_SIZE {
		return "", errors.New("File size unacceptable")
	}
	imgType := v.valImageMime(img.Header.Get("Content-Type"))
	if imgType == "" {
		return "", errors.New("Unsupported file type")
	}

	return imgType, nil
}
