package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model"
	"github.com/nanoLeinz/librarium/repository"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var intServerErr = myerror.InternalServerErr

type BookCopyServiceImpl struct {
	log      *log.Logger
	copyRepo repository.BookCopyRepository
}

func NewBookCopyService(log *log.Logger, copyRepo repository.BookCopyRepository) BookCopyService {
	return &BookCopyServiceImpl{
		log:      log,
		copyRepo: copyRepo,
	}
}

func (s *BookCopyServiceImpl) logWithCtx(ctx context.Context, function string) *logrus.Entry {

	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

	logger := s.log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": function,
	})

	return logger
}

func (s *BookCopyServiceImpl) Create(ctx context.Context, bookId uuid.UUID, status string, copies int) error {

	err := s.copyRepo.Create(ctx, bookId, status, copies)
	if err != nil {
		return intServerErr
	}

	return nil
}
func (s *BookCopyServiceImpl) Update(ctx context.Context, bookCopy *model.BookCopy) error {
	return nil
}
func (s *BookCopyServiceImpl) DeleteById(ctx context.Context, bookCopyId uint) error {
	return nil
}
func (s *BookCopyServiceImpl) GetByID(ctx context.Context, bookCopyId uint) (*model.BookCopy, error) {
	return nil, nil
}
func (s *BookCopyServiceImpl) GetAll(ctx context.Context) (*[]model.BookCopy, error) {
	return nil, nil
}
func (s *BookCopyServiceImpl) GetByCondition(ctx context.Context, bookCopy *model.BookCopy) (*[]model.BookCopy, error) {
	return nil, nil
}
