package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
)

type BookCopyRepository interface {
	Create(ctx context.Context, bookId uuid.UUID, status string, copies int) error
	Update(ctx context.Context, bookCopy *model.BookCopy) error
	DeleteById(ctx context.Context, bookCopyId uint) error
	GetByID(ctx context.Context, bookCopyId uint) (*model.BookCopy, error)
	GetAll(ctx context.Context) (*[]model.BookCopy, error)
}
