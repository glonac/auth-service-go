package storage

import (
	"auth-service/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

var ErrNoMatch = fmt.Errorf("No mathing redord")

func Initialize(cfDb *config.ConfigDatabase) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfDb.Host, cfDb.Port, cfDb.Name, cfDb.Password, cfDb.Name)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{})

	if err != nil {
		log.Println(err)
	}
	log.Println("Database connected")

	DB = db
	return db
}
