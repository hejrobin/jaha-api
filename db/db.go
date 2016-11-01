package db

import (
	// Native packages
	"log"
	"os"
	"sync"

	// 3rd party packages
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	// Local packages
	"jaha-api/utils"
)

var connection *gorm.DB
var initOnce sync.Once

/**
 *	Returns database connection instance, creates a new instance if not set.
 *
 *	@return *gorm.DB
 */
func GetConnection() *gorm.DB {
	var connectionString string
	var connectionError error

	initOnce.Do(func() {
		connectionString = os.Getenv("DSN")
		driverName := utils.Pick(os.Getenv("DBDRIVER"), "mysql")

		if connectionString == "" {
			log.Fatal("DSN not set.")
		}

		connectionString = connectionString + "?charset=utf8&parseTime=True"
		connection, connectionError = gorm.Open(driverName, connectionString)

		if connectionError != nil {
			log.Fatalln(connectionError)
		}
	})

	return connection
}
