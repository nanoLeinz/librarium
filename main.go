package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/nanoLeinz/librarium/controller"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/repository"
	"github.com/nanoLeinz/librarium/router"
	"github.com/nanoLeinz/librarium/service"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	log.SetLevel(log.TraceLevel)
	// log.SetReportCaller(true)
}

func main() {
	log.Info("App Starting")
	err := godotenv.Load()

	if err != nil {
		log.WithError(err).Error("failed loading env")
	}

	db := helper.InitDatabase()
	helper.AutoMigrateModels(db)

	MemberRepo := repository.NewMemberRepository(db, log.StandardLogger())
	MemberServ := service.NewMemberServiceImpl(MemberRepo, log.StandardLogger())

	validate := validator.New()

	MemberHandler := controller.NewMemberController(MemberServ, validate, log.StandardLogger())
	AuthHandler := controller.NewAuthController(MemberServ, validate, log.StandardLogger())

	AuthorRepo := repository.NewAuthorRepositoryImpl(log.StandardLogger(), db)
	AuthorServ := service.NewAuthorServiceImpl(log.StandardLogger(), AuthorRepo)
	AuthorHandler := controller.NewAuthorController(AuthorServ, log.StandardLogger())

	BookCopyRepo := repository.NewBookCopyRepositoryImpl(log.StandardLogger(), db)

	BookRepo := repository.NewBookRepositoryImpl(log.StandardLogger(), db)
	BookServ := service.NewBookServiceImpl(log.StandardLogger(), BookRepo, BookCopyRepo)
	BookHandler := controller.NewBookController(BookServ, log.StandardLogger())

	router := router.NewRouter(MemberHandler, AuthHandler, AuthorHandler, BookHandler)

	server := http.Server{
		Addr:         ":8890",
		Handler:      router,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	log.Infof("Server Started at port %+v\n", ":8890")

	err = server.ListenAndServe()

	if err != nil {
		wrapper := fmt.Errorf("cant start server : %w", err)
		panic(wrapper)
	}

}
