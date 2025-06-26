package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nanoLeinz/librarium/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println(err.Error())
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic("error")

	}

	db.AutoMigrate(&model.Author{})
	db.AutoMigrate(&model.Loan{})
	db.AutoMigrate(&model.BookCopy{})
	db.AutoMigrate(&model.Book{})
	db.AutoMigrate(&model.Fine{})
	db.AutoMigrate(&model.Member{})
	db.AutoMigrate(&model.Reservation{})

}
