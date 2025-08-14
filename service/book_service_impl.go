package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookServiceImpl struct {
	log        *logrus.Logger
	repo       repository.BookRepository
	copyrepo   repository.BookCopyRepository
	authorrepo repository.AuthorRepository
}

func NewBookServiceImpl(log *logrus.Logger, repo repository.BookRepository, copy repository.BookCopyRepository) BookService {
	return &BookServiceImpl{
		log:      log,
		repo:     repo,
		copyrepo: copy,
	}
}

func (s *BookServiceImpl) Create(ctx context.Context, data *dto.BookRequest) (*dto.BookResponse, error) {

	s.log.WithFields(logrus.Fields{
		"function":   "BookService.Create",
		"book title": data.Title,
	}).Info("Starting Insert Book")

	authors := []model.Author{}

	if len(data.AuthorIds) > 0 {
		for i := 0; i < len(data.AuthorIds); i++ {
			author := model.Author{
				Model: gorm.Model{
					ID: uint(data.AuthorIds[i]),
				},
			}

			authors = append(authors, author)

		}
	}

	book := &model.Book{
		Title:           data.Title,
		ISBN:            data.ISBN,
		PublicationYear: data.PublicationYear,
		Genre:           data.Genre,
		Author:          authors,
	}

	s.log.WithFields(logrus.Fields{
		"function": "BookService.Create",
		"book":     book,
	}).Info("Starting Insert Book")

	result, err := s.repo.Create(ctx, book)

	if err != nil {
		s.log.WithError(err).Error("failed adding book")

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr); pgErr.Code == "23505" {
			return nil, myerror.NewDuplicateError("book")
		}

		return nil, myerror.InternalServerErr
	}

	if err := s.copyrepo.Create(ctx, result.ID, "available", int(data.InitialCopy)); err != nil {
		s.log.WithError(err).Error("failed adding book copy")

		return nil, myerror.InternalServerErr
	}

	// s.authorrepo.

	response := dto.ToBookResponse(*result, authors)

	s.log.WithFields(logrus.Fields{
		"function": "BookService.Create",
		"book":     response,
	}).Info("Successfuly Inserted Book")

	return &response, nil

}

// func (s *BookServiceImpl) Update(ctx context.Context, id uuid.UUID, data map[string]any) error
// func (s *BookServiceImpl) DeleteByID(ctx context.Context, id uuid.UUID) error
// func (s *BookServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*dto.BookResponse, error)
// func (s *BookServiceImpl) GetByTitle(ctx context.Context, name string) (*[]dto.BookResponse, error)
// func (s *BookServiceImpl) GetAll(ctx context.Context) (*[]dto.BookResponse, error)
