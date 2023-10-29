package db

import (
	"log"
	"os"

	// "alivassopoli.com/leandro-twin-mais-tema/pkg/models" //TODO: remove this automigration
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	dbURL := os.Getenv("DATABASE_STRING_CONN")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Panicln(err)
	}

	// db.AutoMigrate(&models.Exercise{}) //TODO: remove this automigration

	return db
}
