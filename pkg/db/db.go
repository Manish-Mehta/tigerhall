package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/cbroglie/mustache"
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
	Execute(string, interface{}, bool) (*gorm.DB, error)
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
		})
		errorHandler.CheckErrorAndExit(err, "Error in integrating GO-ORM with the DB")
		client = &dBClient{client: gormDB}
	}
}

var GetDBClient = func(service string) DBClient {
	return client
}

type dBClient struct {
	client *gorm.DB
}

func (db *dBClient) GetClient() *gorm.DB {
	return db.client
}

// TODO: Review if the below function needed?
func (db *dBClient) Execute(queryTemplate string, fields interface{}, isSelectQuery bool) (*gorm.DB, error) {

	query, err := mustache.Render(queryTemplate, fields)

	if err != nil {
		return nil, err
	}

	var tx *gorm.DB
	if isSelectQuery {
		tx = db.client.Raw(query)
	} else {
		tx = db.client.Exec(query)
	}

	err = tx.Error
	return tx, err
}
