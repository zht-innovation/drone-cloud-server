package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB
var sqlDB *sql.DB

func DatabaseSetup() {
	var err error

	config := viper.New()
	config.AddConfigPath("./config")
	config.SetConfigName("db")
	config.SetConfigType("yaml")

	if err = config.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Fatalln("Database config not found.")
		}
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := config.GetString("database.dev.host")
	port := config.GetString("database.dev.port")
	dbname := config.GetString("database.dev.dbname")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host,
		port, dbname)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "zht_",
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatalf("Open database error: %v", err)
	}

	sqlDB, _ = db.DB()
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	sqlDB.SetMaxIdleConns(30)
}

func Close() {
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Close conn pools error: %v", err)
	}
}
