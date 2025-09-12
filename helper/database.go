package helper

import (
	stdlog "log"
	"os"
	"time"

	"github.com/nanoLeinz/librarium/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase() *gorm.DB {

	retry := 5
	waitTime := 10 * time.Second

	newLogger := logger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	for i := 0; i < retry; i++ {

		log.Infof("trying connecting to db : %d attempt", i+1)

		db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
			Logger: newLogger,
		})

		if err == nil {
			log.Info("db connected successfully")

			return db
		}

		log.WithError(err).Errorf("db failed to connect, retrying in %v seconds...", waitTime)

		time.Sleep(waitTime)

	}

	log.Warnf("db failed to connect after %d attemps", retry)

	log.Fatal("Failed to initialize database connection")

	return nil

}

func AutoMigrateModels(db *gorm.DB) {
	db.AutoMigrate(&model.Author{})
	db.AutoMigrate(&model.Loan{})
	db.AutoMigrate(&model.BookCopy{})
	db.AutoMigrate(&model.Book{})
	db.AutoMigrate(&model.Fine{})
	db.AutoMigrate(&model.Member{})
	db.AutoMigrate(&model.Reservation{})
}
