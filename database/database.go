package database

import (
	"api/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func New() {
	DB = initDB()
	createMockUsers()
}

func initDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DB_INFOS")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalln(err)
	}
	database.AutoMigrate(&models.User{})

	return database
}

func createMockUsers() {
	users := []models.User{
		{Nickname: "Hichame", Email: "hichame@42LH.fr", Password: "hichame42LH"},
		{Nickname: "Maxime", Email: "maxime@42LH.fr", Password: "maxime42LH"},
		{Nickname: "Yanis", Email: "yanis@42LH.fr", Password: "yanis42LH"},
		{Nickname: "Omar", Email: "omar@42LH.fr", Password: "omar42LH"},
	}

	for i := range users {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(users[i].Password), 10)
		users[i].Password = string(hashedPassword)
		if err := DB.Create(&users[i]).Error; err != nil {
			log.Printf("Failed to create user %s: %v", users[i].Nickname, err)
		} else {
			log.Printf("Created mock user: %s", users[i].Nickname)
		}
	}
}
