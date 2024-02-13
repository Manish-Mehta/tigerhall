package datastore

import (
	"github.com/Manish-Mehta/tigerhall/model/entities"
	"gorm.io/gorm"
)

type TigerStore interface {
	Create(*entities.Tiger) error
	Update(*entities.Tiger) error
	List(dest *[]*entities.Tiger, page int, limit int) error
	NameExists(string) (bool, error)
}

var NewTigerStore = func(db *gorm.DB) TigerStore {
	return &tigerStore{db: db}
}

type tigerStore struct {
	db *gorm.DB
}

func (ts *tigerStore) NameExists(name string) (bool, error) {
	var count int64
	ts.db.Model(&entities.Tiger{}).Where("name = ?", name).Count(&count)
	return count > 0, nil
}

func (ts *tigerStore) Create(tiger *entities.Tiger) error {
	return ts.db.Create(tiger).Error
}

func (ts *tigerStore) Update(tiger *entities.Tiger) error {
	return ts.db.Model(tiger).Updates(tiger).Error
}

func (ts *tigerStore) Get(dest *entities.Tiger, condition *entities.Tiger, fields []string) error {
	return ts.db.Where(condition).Select(fields).Find(dest).Error
}

func (ts *tigerStore) List(dest *[]*entities.Tiger, page int, limit int) error {
	offset := (page - 1) * limit
	return ts.db.Order("last_seen desc").Limit(limit).Offset(offset).Find(dest).Error
}
