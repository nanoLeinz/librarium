package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
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

func (s *BookServiceImpl) logWithCtx(ctx context.Context, function string) *logrus.Entry {

	traceID := ctx.Value("traceID")

	logger := s.log.WithFields(logrus.Fields{
		"traceID":  traceID,
		"function": function,
	})

	return logger
}

func (s *BookServiceImpl) Create(ctx context.Context, data *dto.BookRequest) (*dto.BookResponse, error) {

	logger := s.logWithCtx(ctx, "BookService.Create")

	logger.WithField("bookTitle", data.Title).Info("executing query")

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

	logger.WithField("book", book).Info()

	result, err := s.repo.Create(ctx, book)

	if err != nil {
		logger.WithError(err).Error("failed executing query")

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr); pgErr.Code == "23505" {
			return nil, myerror.NewDuplicateError("book")
		}

		return nil, myerror.InternalServerErr
	}

	logger.Info("executing insert book copy query")
	if err := s.copyrepo.Create(ctx, result.ID, "available", int(data.InitialCopy)); err != nil {
		logger.WithError(err).Error("failed to execute insert book copy query")
		return nil, myerror.InternalServerErr
	}

	fetchedAuthors, err := s.authorrepo.GetByIDs(ctx, data.AuthorIds...)

	if err != nil {
		logger.WithError(err).Error("failed to execute get author query")
		return nil, myerror.InternalServerErr
	}

	response := dto.ToBookResponse(*result, *fetchedAuthors)
	logger.WithField("response", response).Info("query executed successfully")

	return &response, nil

}

func (s *BookServiceImpl) Update(ctx context.Context, id uuid.UUID, book *dto.BookRequest) error {

	logger := s.logWithCtx(ctx, "BookService.Update").
		WithField("book", book.Title)

	logger.Info("Processing request to update member")

	var data = dto.ToBookModel(id, *book)

	err := s.repo.Update(ctx, &data)

	if err != nil {
		logger.WithError(err).Error("failed to process update request")

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr); pgErr.Code == "23505" {
			return myerror.NewNotFoundError("book")
		} else {
			return myerror.InternalServerErr
		}

	}
	return nil
}
func (s *BookServiceImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {

	return nil
}
func (s *BookServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*dto.BookResponse, error)
func (s *BookServiceImpl) GetByTitle(ctx context.Context, name string) (*[]dto.BookResponse, error)
func (s *BookServiceImpl) GetAll(ctx context.Context) (*[]dto.BookResponse, error)
