package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model/dto"
	log "github.com/sirupsen/logrus"
)

type BookCopyService interface {
	logWithCtx(ctx context.Context, function string) *log.Entry
	Create(ctx context.Context, bookId uuid.UUID, status string, copies int) error
	Update(ctx context.Context, copyId uint, bookCopy *dto.BookCopyRequest) error
	DeleteById(ctx context.Context, bookCopyId uint) error
	GetByID(ctx context.Context, bookCopyId uint) (*dto.BookCopyResponse, error)
	GetAll(ctx context.Context) (*[]dto.BookCopyResponse, error)
	GetByCondition(ctx context.Context, bookCopy *dto.BookCopyRequest) (*[]dto.BookCopyResponse, error)
}
