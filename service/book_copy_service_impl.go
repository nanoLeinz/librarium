package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/repository"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	errIntServer = myerror.InternalServerErr
	errNotFound  = myerror.NewNotFoundError("copy")
)

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

func (s *BookCopyServiceImpl) logWithCtx(ctx context.Context, function string) *log.Entry {

	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

	logger := s.log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": function,
	})

	return logger
}

func (s *BookCopyServiceImpl) Create(ctx context.Context, bookId uuid.UUID, status string, copies int) error {
	logger := s.logWithCtx(ctx, "BookCopyService.Create").
		WithFields(log.Fields{
			"bookId": bookId,
			"status": status,
			"copies": copies,
		})

	logger.Info("received create book copy request")

	err := s.copyRepo.Create(ctx, bookId, status, copies)
	if err != nil {
		logger.WithError(err).Error("failed to create book copies")
		return errIntServer
	}

	logger.Info("book copies created successfully")
	return nil
}

func (s *BookCopyServiceImpl) Update(ctx context.Context, copyId uint, bookCopy *dto.BookCopyRequest) error {
	logger := s.logWithCtx(ctx, "BookCopyService.Update").
		WithFields(log.Fields{
			"copyId": copyId,
			"status": bookCopy.Status,
			"bookId": bookCopy.BookID,
		})

	logger.Info("received update book copy request")

	copy := model.BookCopy{
		Model: gorm.Model{
			ID: copyId,
		},
		Status: bookCopy.Status,
	}

	err := s.copyRepo.Update(ctx, &copy)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithError(err).Error("book copy not found")
			return errNotFound
		}
		logger.WithError(err).Error("failed to update book copy")
		return errIntServer
	}

	logger.Info("book copy updated successfully")
	return nil
}

func (s *BookCopyServiceImpl) DeleteById(ctx context.Context, bookCopyId uint) error {
	logger := s.logWithCtx(ctx, "BookCopyService.DeleteById").
		WithField("copyId", bookCopyId)

	logger.Info("received delete book copy request")

	err := s.copyRepo.DeleteById(ctx, bookCopyId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithError(err).Error("book copy not found")
			return errNotFound
		}
		logger.WithError(err).Error("failed to delete book copy")
		return errIntServer
	}

	logger.Info("book copy deleted successfully")
	return nil
}

func (s *BookCopyServiceImpl) GetByID(ctx context.Context, bookCopyId uint) (*dto.BookCopyResponse, error) {
	logger := s.logWithCtx(ctx, "BookCopyService.GetByID").
		WithField("copyId", bookCopyId)

	logger.Info("received get book copy by ID request")

	rs, err := s.copyRepo.GetByID(ctx, bookCopyId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithError(err).Error("book copy not found")
			return nil, errNotFound
		}
		logger.WithError(err).Error("failed to get book copy by ID")
		return nil, errIntServer
	}

	bookCopyRs := dto.BookCopyResponse{
		ID:     rs.Model.ID,
		Status: rs.Status,
		BookID: rs.BookID,
	}
	logger.Info("book copy fetched successfully")
	return &bookCopyRs, nil
}

func (s *BookCopyServiceImpl) GetAll(ctx context.Context) (*[]dto.BookCopyResponse, error) {
	logger := s.logWithCtx(ctx, "BookCopyService.GetAll")
	logger.Info("received get all book copies request")

	rs, err := s.copyRepo.GetAll(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithError(err).Error("no book copies found")
			return nil, errNotFound
		}
		logger.WithError(err).Error("failed to get all book copies")
		return nil, errIntServer
	}

	bookCopies := []dto.BookCopyResponse{}
	for _, v := range *rs {
		bookCopies = append(bookCopies, dto.BookCopyResponse{
			ID:     v.Model.ID,
			Status: v.Status,
			BookID: v.BookID,
		})
	}

	logger.WithField("count", len(bookCopies)).Info("all book copies fetched successfully")
	return &bookCopies, nil
}

func (s *BookCopyServiceImpl) GetByCondition(ctx context.Context, bookCopy *dto.BookCopyRequest) (*[]dto.BookCopyResponse, error) {
	logger := s.logWithCtx(ctx, "BookCopyService.GetByCondition").
		WithFields(log.Fields{
			"status": bookCopy.Status,
			"bookId": bookCopy.BookID,
		})

	logger.Info("received get book copies by condition request")

	copy := model.BookCopy{
		Status: bookCopy.Status,
		BookID: bookCopy.BookID,
	}

	rs, err := s.copyRepo.GetByCondition(ctx, &copy)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.WithError(err).Error("no book copies found for condition")
			return nil, errNotFound
		}
		logger.WithError(err).Error("failed to get book copies by condition")
		return nil, errIntServer
	}

	bookCopies := []dto.BookCopyResponse{}
	for _, v := range *rs {
		bookCopies = append(bookCopies, dto.BookCopyResponse{
			ID:     v.Model.ID,
			Status: v.Status,
			BookID: v.BookID,
		})
	}

	logger.WithField("count", len(bookCopies)).Info("book copies fetched by condition successfully")
	return &bookCopies, nil
}
