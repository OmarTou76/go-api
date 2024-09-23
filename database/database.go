package database

import (
    "log"
    "os"
    "api/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func New() {
    DB = initDB()
}

func initDB() *gorm.DB  {
    err := godotenv.Load() 
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    database, err := gorm.Open(postgres.Open(os.Getenv("DB_INFOS")), &gorm.Config{})
    if err != nil {
        log.Fatalln(err)
    }
    database.AutoMigrate(&models.User{})
    return database
} 

