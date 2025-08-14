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

func (s *AuthorServiceImpl) Create(ctx context.Context, data *dto.AuthorRequest) (*dto.AuthorResponse, error) {

	s.log.WithFields(log.Fields{
		"function": "AuthorService.Create",
		"author":   *data,
	}).Info("receive create request from controller")

	author := model.Author{
		Name:      data.Name,
		Biography: data.Biography,
		BirthYear: data.BirthYear,
	}

	result, err := s.repo.Create(ctx, &author)

	if err != nil {
		s.log.WithError(err).Error("error inserting author")

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr); pgErr.Code == "23505" {
			return nil, myerror.NewDuplicateError("author")
		} else {
			return nil, myerror.InternalServerErr
		}
	}

	s.log.WithFields(log.Fields{
		"author": *result,
	}).Info("receive response from repository")

	response := dto.ToAuthorResponse(*result)

	s.log.WithFields(log.Fields{
		"data": response,
	}).Info("converted to dto and send response")

	return &response, nil
}
