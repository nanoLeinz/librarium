package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/repository"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	ErrNotFoundAuthor  = myerror.NewNotFoundError("author")
	ErrDuplicateAuthor = myerror.NewDuplicateError("author")
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
	traceID := ctx.Value("traceID")

	return log.WithFields(log.Fields{
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

	logger := s.logWithCtx(ctx, "AuthorService.Create").
		WithFields(log.Fields{
			"authorName": data.Name,
			"authorID":   id,
		})

	logger.Info("receive create request from controller")

	author := model.Author{
		Model:     gorm.Model{ID: id},
		Name:      data.Name,
		Biography: data.Biography,
		BirthYear: data.BirthYear,
	}

	err := s.repo.Update(ctx, author)

	if err != nil {
		logger.WithError(err).Error("failed to update author to repository")

		var e *pgconn.PgError
		if errors.As(err, &e) {
			switch e.Code {
			case "23505":
				logger.WithError(myerror.NewNotFoundError("author")).Error("error converted to notfound")
				return myerror.NewNotFoundError("author")
			default:
				logger.WithError(myerror.InternalServerErr).Error("error converted to internalservererror")
				return myerror.InternalServerErr
			}
		}
	}

	return nil

}

func (s *AuthorServiceImpl) DeleteById(ctx context.Context, id uint) error
func (s *AuthorServiceImpl) GetByIDs(ctx context.Context, ids ...uint) (*[]dto.AuthorResponse, error)
func (s *AuthorServiceImpl) GetAll(ctx context.Context) (*[]dto.AuthorResponse, error)
func (s *AuthorServiceImpl) GetAuthorsBook(ctx context.Context, author *dto.AuthorRequest) (*[]dto.AuthorResponse, error)
