package datastore

import (
	"github.com/Manish-Mehta/tigerhall/model/entities"
	"gorm.io/gorm"
)

type UserStore interface {
	Create(*entities.User) error
	Get(dest *entities.User, condition *entities.User, fields []string) error
	EmailExists(string) (bool, error)
}

var NewUserStore = func(db *gorm.DB) UserStore {
	return &userStore{db: db}
}

type userStore struct {
	db *gorm.DB
}

func (us *userStore) EmailExists(email string) (bool, error) {
	var count int64
	us.db.Model(&entities.User{}).Where("email = ?", email).Count(&count)
	return count > 0, nil
}

func (us *userStore) Create(user *entities.User) error {
	return us.db.Create(user).Error
}

func (us *userStore) Get(dest *entities.User, condition *entities.User, fields []string) error {
	return us.db.Where(condition).Select(fields).Find(dest).Error
}
