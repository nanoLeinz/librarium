package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookCopyRepositoryImpl struct {
	log *log.Logger
	db  *gorm.DB
}

func NewBookCopyRepositoryImpl(log *log.Logger, db *gorm.DB) BookCopyRepository {
	return &BookCopyRepositoryImpl{
		log: log,
		db:  db,
	}
}

func (s *BookCopyRepositoryImpl) logWithCtx(ctx context.Context, function string) *log.Entry {

	traceId := ctx.Value("traceID")

	logging := s.log.WithFields(log.Fields{
		"function": function,
		"traceID":  traceId,
	})

	return logging
}

func (s *BookCopyRepositoryImpl) Create(ctx context.Context, bookId uuid.UUID, status string, copies int) error {

	logger := s.logWithCtx(ctx, "BookRepository.Create")

	logger.WithFields(log.Fields{
		"bookId": bookId,
		"copies": copies,
	}).Info("executing insert book copy query")

	bookCopies := []model.BookCopy{}

	for i := 0; i < copies; i++ {
		var bookCopy = model.BookCopy{
			BookID: bookId,
			Status: status,
		}

		bookCopies = append(bookCopies, bookCopy)
	}

	result := s.db.WithContext(ctx).Create(&bookCopies)

	if result.Error != nil {
		logger.WithError(result.Error).Error("failed executing book copy query")
		return result.Error
	} else if result.RowsAffected == 0 {
		logger.Debug("query executed but 0 data inserted")
	} else {
		logger.WithFields(log.Fields{
			"copies inserted": result.RowsAffected,
		}).Info("successfully executed book copy insert query")
	}
	return nil
}

func (s *BookCopyRepositoryImpl) Update(ctx context.Context, bookCopy *model.BookCopy) error {

	logger := s.logWithCtx(ctx, "BookCopyRepository.Update").WithFields(log.Fields{
		"bookCopyID": bookCopy.ID,
	})

	logger.Info("executing update query")

	result := s.db.WithContext(ctx).Updates(bookCopy)

	if result.Error != nil {
		logger.WithError(result.Error).Error("failed executing update book copy query")

		return result.Error
	} else if result.RowsAffected == 0 {
		logger.Debug("query executed but 0 rows affected")
	} else {
		logger.Info("book copy update query executed successfully")
	}

	return nil
}

func (s *BookCopyRepositoryImpl) DeleteById(ctx context.Context, bookCopyId uint) error {

	logger := s.logWithCtx(ctx, "BookCopyRepository.DeleteById").WithFields(log.Fields{
		"bookCopyID": bookCopyId,
	})

	logger.Info("executing book copy delete query")

	result := s.db.WithContext(ctx).Delete(&model.BookCopy{}, bookCopyId)

	if result.Error != nil {
		logger.WithError(result.Error).Error("failed to execute delete book copy by id")
		return result.Error
	} else if result.RowsAffected == 0 {
		logger.Debug("delete book copy query executed but 0 rows affected")
	} else {
		logger.Info("delete book copy query executed successfully")
	}

	return nil
}

func (s *BookCopyRepositoryImpl) GetByID(ctx context.Context, bookCopyId uint) (*model.BookCopy, error) {

	logger := s.logWithCtx(ctx, "BookCopyRepository.GetByID").WithFields(log.Fields{
		"bookCopyID": bookCopyId,
	})

	logger.Info("executing get by id query")

	bookCopy := &model.BookCopy{}

	err := s.db.WithContext(ctx).First(bookCopy, bookCopyId).Error
	if err != nil {
		logger.WithError(err).Error("failed executing get by id quert")

		return nil, err
	}

	logger.Info("get by id query executed successfully")

	return bookCopy, nil
}

func (s *BookCopyRepositoryImpl) GetAll(ctx context.Context) (*[]model.BookCopy, error) {

	s.logWithCtx(ctx, "BookCopyRepository.GetAll").Info("executing get all query")

	var copies []model.BookCopy

	err := s.db.WithContext(ctx).Scopes(helper.Paginator(ctx)).Find(copies).Error

	if err != nil {
		s.logWithCtx(ctx, "BookCopyRepository.GetAll").Error("failed executing get all query")

		return nil, err
	}

	s.logWithCtx(ctx, "BookCopyRepository.GetAll").Info("get all query executed successfully")
	return nil, nil
}

func (s *BookCopyRepositoryImpl) GetByCondition(ctx context.Context, bookCopy *model.BookCopy) (*[]model.BookCopy, error) {

	logger := s.logWithCtx(ctx, "BookCopyRepository.GetByCondition").WithFields(log.Fields{
		"status": bookCopy.Status,
	})

	logger.Info("executing get by condition query")

	var copies []model.BookCopy

	err := s.db.WithContext(ctx).Scopes(helper.Paginator(ctx)).Where(bookCopy).Find(&copies).Error

	if err != nil {
		logger.WithError(err).Error("failed executing get by condition query")
		return nil, err
	}

	logger.WithField("count", len(copies)).Info("get by condition query executed successfully")
	return &copies, nil

}
