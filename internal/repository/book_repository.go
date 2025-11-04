package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/internal/model"
)

type BookRepository interface {
	Create(ctx context.Context, data *model.Book) (*model.Book, error)
	Update(ctx context.Context, data *model.Book) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Book, error)
	GetByTitle(ctx context.Context, name string) (*[]model.Book, error)
	GetAll(ctx context.Context) (*[]model.Book, error)
	GetBooksAuthor(context.Context, *model.Book) (*[]model.Author, error)
}
