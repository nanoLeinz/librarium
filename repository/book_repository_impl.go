package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookRepositoryImpl struct {
	log *logrus.Logger
	db  *gorm.DB
}

func NewBookRepositoryImpl(log *logrus.Logger, db *gorm.DB) BookRepository {
	return &BookRepositoryImpl{
		log: log,
		db:  db,
	}
}

func (s *BookRepositoryImpl) Create(ctx context.Context, data *model.Book) (*model.Book, error) {

	s.log.WithField("Book Title", data.Title).Info("Inserting to DB")

	if err := s.db.WithContext(ctx).Create(data).Error; err != nil {

		s.log.WithError(err).Error("Error Inserting Book to DB")

		return nil, err
	}

	s.log.WithFields(logrus.Fields{
		"Book ID": data.ID,
		"Title":   data.Title,
	}).Info("Successfully inserted book")

	return data, nil
}

func (s *BookRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data map[string]any) error {

	s.log.WithFields(logrus.Fields{
		"function": "Update",
		"Book ID":  id.String(),
		"data":     data,
	}).Info("Updating to DB")

	if err := s.db.WithContext(ctx).Model(&model.Book{}).Where("id = ?", id).Updates(data).Error; err != nil {

		s.log.WithError(err).Error("failed to update record")

		return err
	}
	s.log.WithFields(logrus.Fields{
		"function": "Update",
		"Book ID":  id.String(),
		"data":     data,
	}).Info("Successfully Updated Book to DB")
	return nil
}

func (s *BookRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	s.log.WithFields(logrus.Fields{
		"function": "DeleteByID",
		"Book ID":  id.String(),
	}).Info("deleting record DB")

	if err := s.db.WithContext(ctx).Delete(&model.Book{}, id).Error; err != nil {

		s.log.WithError(err).Error("failed deleting record")
		return err

	}

	s.log.WithFields(logrus.Fields{
		"function": "DeleteByID",
		"Book ID":  id.String(),
	}).Info("Successfully Deleted Book DB")
	return nil
}

func (s *BookRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Book, error) {

	s.log.WithFields(logrus.Fields{
		"function": "GetByID",
		"Book ID":  id.String(),
	}).Info("fetching record from DB")

	var data = &model.Book{}

	if err := s.db.WithContext(ctx).First(data, id).Error; err != nil {
		s.log.WithError(err).Error("failed fetching record")
		return nil, err

	}

	s.log.WithFields(logrus.Fields{
		"function": "GetByID",
		"Book ID":  id.String(),
	}).Info("Successfully Fetched Book from DB")

	return data, nil
}
func (s *BookRepositoryImpl) GetByTitle(ctx context.Context, name string) (*[]model.Book, error) {

	s.log.WithFields(logrus.Fields{
		"function": "GetByName",
		"name":     name,
	}).Info("fetching record from DB")

	var datas = &[]model.Book{}

	if err := s.db.WithContext(ctx).Where("title LIKE '%?%'", name).Find(datas).Error; err != nil {
		s.log.WithError(err).Error("failed fetching record")

		return nil, err

	}

	s.log.WithFields(logrus.Fields{
		"function": "GetByName",
		"books":    *datas,
	}).Info("Successfully Fetched Book from DB")

	return datas, nil
}
func (s *BookRepositoryImpl) GetAll(ctx context.Context) (*[]model.Book, error) {
	s.log.WithFields(logrus.Fields{
		"function": "GetAll",
	}).Info("fetching record from DB")

	var datas = &[]model.Book{}

	if err := s.db.WithContext(ctx).Find(datas).Error; err != nil {
		s.log.WithError(err).Error("failed fetching record")

		return nil, err

	}

	s.log.WithFields(logrus.Fields{
		"function": "GetAll",
		"books":    *datas,
	}).Info("Successfully Fetched Book from DB")

	return datas, nil
}
