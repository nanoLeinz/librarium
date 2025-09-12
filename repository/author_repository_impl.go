package repository

import (
	"context"

	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/model"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthorRepositoryImpl struct {
	log *log.Logger
	db  *gorm.DB
}

func NewAuthorRepositoryImpl(log *log.Logger, db *gorm.DB) AuthorRepository {
	return &AuthorRepositoryImpl{
		log: log,
		db:  db,
	}
}

func (s *AuthorRepositoryImpl) logWithCtx(ctx context.Context, function string) *log.Entry {

	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

	logger := s.log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": function,
	})

	return logger

}

func (s *AuthorRepositoryImpl) Create(ctx context.Context, author *model.Author) (*model.Author, error) {

	logger := s.logWithCtx(ctx, "AuthorRepository.Create").
		WithField("authorName", author.Name)

	logger.Info("executing query")

	result := s.db.Create(author)

	err := result.Error
	if err != nil {
		logger.WithError(err).Error("failed inserting author")
		return nil, err

	}

	logger.WithField("authorID", author.ID).Info("inserted author")

	return author, nil
}

func (s *AuthorRepositoryImpl) GetByIDs(ctx context.Context, ids ...uint) (*[]model.Author, error) {

	logger := s.logWithCtx(ctx, "AuthorRepository.GetByIDs").
		WithField("id(s)", ids)

	logger.Info("executing query")

	authors := []model.Author{}

	result := s.db.WithContext(ctx).Find(authors, ids)

	if result.Error != nil {
		logger.WithError(result.Error).Error("failed executing query")
		return nil, result.Error
	} else if len(authors) == 0 {
		logger.Debug("no record fetched")
		return nil, nil
	}

	logger.WithField("data", authors).
		Info("successfully executing query")

	return &authors, nil
}

func (s *AuthorRepositoryImpl) Update(ctx context.Context, author model.Author) error {

	logger := s.logWithCtx(ctx, "AuthorRepository.Update").
		WithField("authorID", author.Model.ID)

	logger.Info("executing query")

	result := s.db.WithContext(ctx).Updates(author)

	if result.Error != nil {
		logger.WithError(result.Error).
			Error("failed executing query")
		return result.Error
	} else if result.RowsAffected == 0 {
		logger.Debug("query executed but 0 rows affected")
	} else {
		logger.Info("successfully executed query")
	}

	return nil
}

func (s *AuthorRepositoryImpl) DeleteById(ctx context.Context, id uint) error {
	logger := s.logWithCtx(ctx, "AuthorRepository.DeleteById").
		WithField("authorID", id)

	logger.Info("executing query")

	if err := s.db.WithContext(ctx).Delete(&model.Author{}, id).Error; err != nil {
		logger.WithError(err).Error("failed executing query")
		return err
	}

	logger.Info("successfully executed query")

	return nil
}

func (s *AuthorRepositoryImpl) GetAll(ctx context.Context) (*[]model.Author, error) {
	logger := s.logWithCtx(ctx, "AuthorRepository.GetAll")

	logger.Info("executing query")

	authors := []model.Author{}

	if err := s.db.WithContext(ctx).Scopes(helper.Paginator(ctx)).Find(&authors).Error; err != nil {
		logger.WithError(err).Error("failed executing query")
		return nil, err
	} else if len(authors) == 0 {
		logger.Debug("no records of author found")
		return nil, gorm.ErrRecordNotFound
	} else {
		logger.Info("query executed succesfully")
	}

	return &authors, nil
}

func (s *AuthorRepositoryImpl) GetAuthorsBook(ctx context.Context, author *model.Author) (*[]model.Book, error) {

	logger := s.logWithCtx(ctx, "AuthorRepository.GetAuthorsBook").
		WithField("authorName", author.Name)

	logger.Info("executing query")

	books := []model.Book{}

	err := s.db.WithContext(ctx).Model(author).Association("Book").Find(&books)

	if err != nil {
		logger.WithError(err).Error("failed to execute query")
		return nil, err
	}

	logger.Info("query executed successfully")
	return &books, nil
}
