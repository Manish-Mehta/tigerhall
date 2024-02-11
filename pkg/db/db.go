package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq" // for postgres driver init
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/Manish-Mehta/tigerhall/internal/config"
	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
)

var client DBClient = nil

type DBClient interface {
	GetClient() *gorm.DB
}

var InitService = func() {
	if client == nil {
		defer log.Println("DB Connected Successfully")

		db, err := sql.Open("postgres", config.DB_STR)
		errorHandler.CheckErrorAndExit(err, "Error in connecting to the DB")
		db.SetConnMaxIdleTime(5 * time.Minute)

		err = db.Ping()
		errorHandler.CheckErrorAndExit(err, "Init DB Ping failed for the DB")

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
			/* ADD ADDITIONAL DB CONFIG HERE */
		}), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			// Logger: logger.Default.LogMode(logger.Info),
		})
		errorHandler.CheckErrorAndExit(err, "Error in integrating GO-ORM with the DB")
		client = &dBClient{client: gormDB}
	}
}

var GetDBClient = func() DBClient {
	return client
}

type dBClient struct {
	client *gorm.DB
}

func (db *dBClient) GetClient() *gorm.DB {
	return db.client
}
