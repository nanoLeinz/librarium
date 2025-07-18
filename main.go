package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/nanoLeinz/librarium/controller"
	"github.com/nanoLeinz/librarium/model"
	"github.com/nanoLeinz/librarium/repository"
	"github.com/nanoLeinz/librarium/router"
	"github.com/nanoLeinz/librarium/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

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

	MemberRepo := repository.NewMemberRepository(db)

	MemberServ := service.NewMemberServiceImpl(MemberRepo)

	validate := validator.New()

	MemberHandler := controller.NewMemberController(MemberServ, validate)

	router := router.NewRouter(MemberHandler)

	server := http.Server{
		Addr:         ":8890",
		Handler:      router,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	log.Printf("Server Started at port %+v\n", ":8890")

	err = server.ListenAndServe()

	if err != nil {
		wrapper := fmt.Errorf("cant start server %w", err)
		panic(wrapper)
	}

}
