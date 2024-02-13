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
	List(dest *[]*entities.Sight, page, limit int, fields []string) error
}

var NewSightStore = func(db *gorm.DB) SightStore {
	return &sightStore{db: db}
}

type sightStore struct {
	db *gorm.DB
}

func (ss *sightStore) NameExists(name string) (bool, error) {
	var count int64
	ss.db.Model(&entities.Sight{}).Where("name = ?", name).Count(&count)
	return count > 0, nil
}

func (ss *sightStore) Create(user *entities.Sight) error {
	return ss.db.Create(user).Error
}

func (ss *sightStore) GetLatest(dest *entities.Sight, condition *entities.Sight, fields []string) error {
	return ss.db.Select(fields).Order("seen_at desc").Find(dest, condition).Limit(1).Error
}

func (ss *sightStore) GetDistance(coord1, coord2 dto.Coordinate) (float64, error) {
	var distance float64

	q := fmt.Sprintf(
		`SELECT ST_Distance('SRID=4326;POINT(%.8f %.8f)'::geometry, 'SRID=4326;POINT(%.8f %.8f)'::geometry) * 40075 / 360 AS distance_km`,
		coord1.Lat, coord1.Lon,
		coord2.Lat, coord2.Lon,
	)
	err := ss.db.Raw(q).Row().Scan(&distance)

	log.Println(distance)
	return distance, err
}

func (ss *sightStore) List(dest *[]*entities.Sight, page, limit int, fields []string) error {
	offset := (page - 1) * limit
	return ss.db.Select(fields).Order("seen_at desc").Limit(limit).Offset(offset).Find(dest).Error
}
