package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
)

type BookService interface {
	Create(ctx context.Context, data *model.Book) (*model.Book, error)
	Update(ctx context.Context, id uuid.UUID, data map[string]any) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Book, error)
	GetByTitle(ctx context.Context, name string) (*[]model.Book, error)
	GetAll(ctx context.Context) (*[]model.Book, error)
}
