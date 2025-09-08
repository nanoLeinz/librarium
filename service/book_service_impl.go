package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	ErrNotFound  = myerror.NewNotFoundError("book")
	ErrDuplicate = myerror.NewDuplicateError("book")
	ErrIntServer = myerror.InternalServerErr
)

type BookServiceImpl struct {
	log      *logrus.Logger
	repo     repository.BookRepository
	copyrepo repository.BookCopyRepository
}

func NewBookServiceImpl(log *logrus.Logger, repo repository.BookRepository, copy repository.BookCopyRepository) BookService {
	return &BookServiceImpl{
		log:      log,
		repo:     repo,
		copyrepo: copy,
	}
}

func (s *BookServiceImpl) logWithCtx(ctx context.Context, function string) *logrus.Entry {

	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

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

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicate
		}

		return nil, ErrIntServer
	}

	logger.Info("executing insert book copy query")
	if err := s.copyrepo.Create(ctx, result.ID, "available", int(data.InitialCopy)); err != nil {
		logger.WithError(err).Error("failed to execute insert book copy query")
		return nil, ErrIntServer
	}

	fetchedAuthors, err := s.repo.GetBooksAuthor(ctx, book)

	if err != nil {
		logger.WithError(err).Error("failed to execute get author query")
		return nil, ErrIntServer
	}

	result.Author = *fetchedAuthors

	// fetchedAuthors, err := s.authorrepo.GetByIDs(ctx, data.AuthorIds...)

	// if err != nil {
	// 	logger.WithError(err).Error("failed to execute get author query")
	// 	return nil, ErrIntServer
	// }

	response := dto.ToBookResponse(*result)

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

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicate
		} else if gorm.ErrRecordNotFound == err {
			return ErrNotFound
		} else {
			return ErrIntServer
		}

	}
	return nil
}
func (s *BookServiceImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {

	logger := s.logWithCtx(ctx, "BookService.DeleteByID").
		WithField("bookID", id)

	logger.Info("received request to delete book")

	err := s.repo.DeleteByID(ctx, id)

	if err != nil {
		logger.WithError(err).Error("failed to delete book from repository")

		if gorm.ErrRecordNotFound == err {
			return ErrNotFound
		} else {
			return ErrIntServer
		}
	}

	logger.Info("successfully deleted book from repository")
	return nil
}

func (s *BookServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*dto.BookResponse, error) {
	logger := s.logWithCtx(ctx, "BookService.GetByID").
		WithField("bookID", id)

	logger.Info("received request to fetch book by id")

	book, err := s.repo.GetByID(ctx, id)

	if err != nil {
		logger.WithError(err).Error("failed to fetch book from repository")

		if gorm.ErrRecordNotFound == err {
			return nil, ErrNotFound
		} else {
			return nil, ErrIntServer
		}
	}

	logger.Info("successfully fetched book from repository")
	logger.Info("fetching author")

	fetchedAuthors, err := s.repo.GetBooksAuthor(ctx, book)

	if err != nil {
		logger.WithError(err).Error("failed to execute get author query")
		return nil, ErrIntServer
	}

	book.Author = *fetchedAuthors

	bookResponse := dto.ToBookResponse(*book)

	return &bookResponse, nil
}
func (s *BookServiceImpl) GetByTitle(ctx context.Context, name string) (*[]dto.BookResponse, error) {

	logger := s.logWithCtx(ctx, "BookService.GetByTitle").
		WithField("bookName", name)

	logger.Info("received request to fetch book by name")

	result, err := s.repo.GetByTitle(ctx, name)

	if err != nil {
		logger.WithError(err).Error("failed to fetch books from repo")

		if gorm.ErrRecordNotFound == err {
			return nil, ErrNotFound
		} else {
			return nil, ErrIntServer
		}
	}

	response := []dto.BookResponse{}

	for _, v := range *result {
		response = append(response, dto.ToBookResponse(v))
	}

	logger.Info("successfully fetched books")
	return &response, nil
}

func (s *BookServiceImpl) GetAll(ctx context.Context) (*[]dto.BookResponse, error) {

	logger := s.logWithCtx(ctx, "BookService.GetAll")
	logger.Info("received request to fetch all books")

	result, err := s.repo.GetAll(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to fetch books from repo")

		if gorm.ErrRecordNotFound == err {
			return nil, ErrNotFound
		} else {
			return nil, ErrIntServer
		}
	}

	response := []dto.BookResponse{}

	for _, v := range *result {
		response = append(response, dto.ToBookResponse(v))
	}

	logger.Info("successfully fetched books")
	return &response, nil

}
