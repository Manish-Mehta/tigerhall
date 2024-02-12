package datastore

import (
	"github.com/Manish-Mehta/tigerhall/model/entities"
	"gorm.io/gorm"
)

type TigerStore interface {
	Create(*entities.Tiger) error
	Get(dest *entities.Tiger, condition *entities.Tiger, fields []string) error
	NameExists(string) (bool, error)
}

var NewTigerStore = func(db *gorm.DB) TigerStore {
	return &tigerStore{db: db}
}

type tigerStore struct {
	db *gorm.DB
}

func (us *tigerStore) NameExists(name string) (bool, error) {
	var count int64
	us.db.Model(&entities.Tiger{}).Where("name = ?", name).Count(&count)
	return count > 0, nil
}

func (us *tigerStore) Create(user *entities.Tiger) error {
	return us.db.Create(user).Error
}

func (us *tigerStore) Get(dest *entities.Tiger, condition *entities.Tiger, fields []string) error {
	return us.db.Where(condition).Select(fields).Find(dest).Error
}
