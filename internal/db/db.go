package db

import (
	"fmt"
	"log"
	"os"

	"rio/internal/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/mysql"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

func ConnectDataBase() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	Dbdriver := os.Getenv("DB_DRIVER")
	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	DB, err = gorm.Open(Dbdriver, DBURL)

	if err != nil {
		fmt.Println("Cannot connect to database", Dbdriver)
		log.Fatal("connection error: ", err)
	} else {
		fmt.Println("We are connected to the database ", Dbdriver)
	}

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Server{})
	DB.AutoMigrate(&models.UserServer{})
}
