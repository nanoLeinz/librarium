package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/internal/model/dto"
)

type BookService interface {
	Create(ctx context.Context, data *dto.BookRequest) (*dto.BookResponse, error)
	Update(ctx context.Context, id uuid.UUID, data *dto.BookRequest) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.BookResponse, error)
	GetByTitle(ctx context.Context, name string) (*[]dto.BookResponse, error)
	GetAll(ctx context.Context) (*[]dto.BookResponse, error)
}
