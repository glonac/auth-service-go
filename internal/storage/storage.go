package storage

import (
	"auth-service/internal/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// зачем здесь глобальная переменная?
var DB *gorm.DB

// errors.New
var ErrNoMatch = fmt.Errorf("No mathing redord")

func Initialize(cfDb *config.ConfigDatabase) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfDb.Host, cfDb.Port, cfDb.Name, cfDb.Password, cfDb.Name)

	// не передаете логгер, gorm будет логировать своим, что плохо
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{})

	if err != nil {
		log.Println(err)
		// return??
	}
	log.Println("Database connected")

	DB = db
	return db
}
