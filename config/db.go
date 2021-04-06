package config

import (
	"apinews/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error to load Env")
	}
	user := os.Getenv("DB_User")
	host := os.Getenv("DB_Host")
	port := os.Getenv("DB_Port")
	name := os.Getenv("DB_Name")
	pass := os.Getenv("DB_Pass")
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)
	db, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{})
	if err != nil {
		fmt.Println("Error", err)
		panic(err)
	}
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Chapter{}, &models.Url{})
	fmt.Println("Database Connected ma brader")
	return db
}
