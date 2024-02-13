package datastore

import (
	"fmt"
	"log"

	"github.com/Manish-Mehta/tigerhall/dto"
	"github.com/Manish-Mehta/tigerhall/model/entities"
	"gorm.io/gorm"
)

type SightStore interface {
	Create(*entities.Sight) error
	GetLatest(dest *entities.Sight, condition *entities.Sight, fields []string) error
	GetDistance(dto.Coordinate, dto.Coordinate) (float64, error)
}

var NewSightStore = func(db *gorm.DB) SightStore {
	return &sightStore{db: db}
}

type sightStore struct {
	db *gorm.DB
}

func (ts *sightStore) NameExists(name string) (bool, error) {
	var count int64
	ts.db.Model(&entities.Sight{}).Where("name = ?", name).Count(&count)
	return count > 0, nil
}

func (ts *sightStore) Create(user *entities.Sight) error {
	return ts.db.Create(user).Error
}

func (ts *sightStore) GetLatest(dest *entities.Sight, condition *entities.Sight, fields []string) error {
	return ts.db.Select(fields).Order("seen_at desc").Find(dest, condition).Limit(1).Error
}

func (ts *sightStore) GetDistance(coord1, coord2 dto.Coordinate) (float64, error) {
	var distance float64

	q := fmt.Sprintf(
		`SELECT ST_Distance('SRID=4326;POINT(%.8f %.8f)'::geometry, 'SRID=4326;POINT(%.8f %.8f)'::geometry) * 40075 / 360 AS distance_km`,
		coord1.Lat, coord1.Lon,
		coord2.Lat, coord2.Lon,
	)
	err := ts.db.Raw(q).Row().Scan(&distance)

	log.Println(distance)
	return distance, err
}
