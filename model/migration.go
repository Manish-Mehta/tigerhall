package model

import (
	"log"

	"github.com/Manish-Mehta/tigerhall/model/entities"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("Migrating DB")
	tx := db.Begin()

	// Execute the raw SQL query to create the postgis extension for Tiger Location
	if err := tx.Exec("CREATE EXTENSION IF NOT EXISTS postgis;").Error; err != nil {
		tx.Rollback()
		return err
	}

	// User table
	if err := tx.AutoMigrate(&entities.User{}); err != nil {
		tx.Rollback()
		return err
	}

	// Tiger table
	if err := tx.AutoMigrate(&entities.Tiger{}); err != nil {
		tx.Rollback()
		return err
	}

	// Sighting table
	if err := tx.AutoMigrate(&entities.Sight{}); err != nil {
		tx.Rollback()
		return err
	}

	//
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
