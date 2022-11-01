package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	err error
)

func newDB(dbDSN string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// DbConnection create database connection
func DbConnection(dbDSN string) error {
	var db = DB

	logMode := viper.GetBool("DB_LOG_MODE")

	loglevel := logger.Silent
	if logMode {
		loglevel = logger.Info
	}

	dbConn, err := newDB(dbDSN)
	if err != nil {
		log.Fatalf("DB connection error: %v", err.Error())
	}

	db, err = gorm.Open(mysql.New(mysql.Config{Conn: dbConn}), &gorm.Config{
		Logger: logger.Default.LogMode(loglevel),
	})
	if err != nil {
		log.Fatalf("Gorm connection error: %v", err.Error())
		return err
	}
	DB = db
	return nil
}

// GetDB connection
func GetDB() *gorm.DB {
	return DB
}
