package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nanoLeinz/librarium/internal/helper"
	"github.com/nanoLeinz/librarium/internal/model"
	"github.com/nanoLeinz/librarium/internal/model/dto"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/internal/repository"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	ErrNotFoundAuthor  = myerror.NewNotFoundError("author")
	ErrDuplicateAuthor = myerror.NewDuplicateError("author")
	ErrIntServerAuthor = myerror.InternalServerErr
)

type AuthorServiceImpl struct {
	log  *log.Logger
	repo repository.AuthorRepository
}

func NewAuthorServiceImpl(log *log.Logger, repo repository.AuthorRepository) AuthorService {
	return &AuthorServiceImpl{
		log:  log,
		repo: repo,
	}
}

func (s *AuthorServiceImpl) logWithCtx(ctx context.Context, function string) *log.Entry {
	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

	return s.log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": function,
	})
}

func (s *AuthorServiceImpl) Create(ctx context.Context, data *dto.AuthorRequest) (*dto.AuthorResponse, error) {

	logger := s.logWithCtx(ctx, "AuthorService.Create").
		WithField("authorName", data.Name)

	logger.Info("receive create request from handler")

	author := model.Author{
		Name:      data.Name,
		Biography: data.Biography,
		BirthYear: data.BirthYear,
	}

	result, err := s.repo.Create(ctx, &author)

	if err != nil {
		logger.WithError(err).Error("failed to create author for repository")

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, myerror.NewDuplicateError("author")
		} else {
			return nil, myerror.InternalServerErr
		}
	}

	logger.WithField("authorID", result.ID).Info("receive response from repository")

	response := dto.ToAuthorResponse(*result)

	logger.Info("converted to dto and send response to handler")

	return &response, nil
}

func (s *AuthorServiceImpl) Update(ctx context.Context, id uint, data *dto.AuthorRequest) error {
	logger := s.logWithCtx(ctx, "AuthorService.Update").
		WithFields(log.Fields{
			"authorName": data.Name,
			"authorID":   id,
		})

	logger.Info("receive update request from handler")

	author := model.Author{
		Model:     gorm.Model{ID: id},
		Name:      data.Name,
		Biography: data.Biography,
		BirthYear: data.BirthYear,
	}

	err := s.repo.Update(ctx, author)
	if err != nil {
		logger.WithError(err).Error("failed to update author in repository")

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			logger.WithError(myerror.NewDuplicateError("author")).Error("duplicate author error")
			return myerror.NewDuplicateError("author")
		} else if err == gorm.ErrRecordNotFound {
			logger.WithError(myerror.NewNotFoundError("author")).Error("author not found")
			return myerror.NewNotFoundError("author")
		} else {
			logger.WithError(myerror.InternalServerErr).Error("internal server error")
			return myerror.InternalServerErr
		}
	}

	logger.WithField("authorID", id).Info("author updated successfully")
	return nil
}

func (s *AuthorServiceImpl) DeleteById(ctx context.Context, id uint) error {
	logger := s.logWithCtx(ctx, "AuthorService.DeleteById").
		WithField("authorID", id)

	logger.Info("receive delete request from handler")

	err := s.repo.DeleteById(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithError(ErrNotFoundAuthor).Error("author not found")
			return ErrNotFoundAuthor
		} else {
			logger.WithError(ErrIntServerAuthor).Error("internal server error on delete")
			return ErrIntServerAuthor
		}
	}

	logger.WithField("authorID", id).Info("author deleted successfully")
	return nil
}

func (s *AuthorServiceImpl) GetByIDs(ctx context.Context, ids ...uint) (*[]dto.AuthorResponse, error) {
	logger := s.logWithCtx(ctx, "AuthorService.GetByIDs").
		WithField("authorIDs", ids)

	logger.Info("receive get by IDs request from handler")

	result, err := s.repo.GetByIDs(ctx, ids...)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithError(ErrNotFoundAuthor).Error("authors not found")
			return nil, ErrNotFoundAuthor
		} else {
			logger.WithError(ErrIntServerAuthor).Error("internal server error on get by IDs")
			return nil, ErrIntServerAuthor
		}
	}

	var response []dto.AuthorResponse
	for _, v := range *result {
		response = append(response, dto.ToAuthorResponse(v))
	}

	logger.WithField("count", len(response)).Info("authors fetched successfully")
	return &response, nil
}

func (s *AuthorServiceImpl) GetAll(ctx context.Context) (*[]dto.AuthorResponse, error) {
	logger := s.logWithCtx(ctx, "AuthorService.GetAll")
	logger.Info("receive get all request from handler")

	result, err := s.repo.GetAll(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithError(ErrNotFoundAuthor).Error("no authors found")
			return nil, ErrNotFoundAuthor
		} else {
			logger.WithError(ErrIntServerAuthor).Error("internal server error on get all")
			return nil, ErrIntServerAuthor
		}
	}

	var response []dto.AuthorResponse
	for _, v := range *result {
		response = append(response, dto.ToAuthorResponse(v))
	}

	logger.WithField("count", len(response)).Info("all authors fetched successfully")
	return &response, nil
}

func (s *AuthorServiceImpl) GetAuthorsBook(ctx context.Context, authorReq *dto.AuthorRequest) (*[]dto.BookResponseShort, error) {
	logger := s.logWithCtx(ctx, "AuthorService.GetAuthorsBook").
		WithField("authorID", authorReq.ID)

	logger.Info("receive get author's books request from handler")

	author := model.Author{
		Model: gorm.Model{
			ID: authorReq.ID,
		},
	}

	result, err := s.repo.GetAuthorsBook(ctx, &author)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithError(ErrNotFound).Error("author's books not found")
			return nil, ErrNotFound
		} else {
			logger.WithError(ErrIntServer).Error("internal server error on get author's books")
			return nil, ErrIntServer
		}
	}

	var response []dto.BookResponseShort
	for _, v := range *result {
		response = append(response, dto.ToBookResponseShort(v))
	}

	logger.WithField("count", len(response)).Info("author's books fetched successfully")
	return &response, nil
}
