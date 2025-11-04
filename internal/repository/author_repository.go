package repository

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/nanoLeinz/librarium/internal/model"
)

type AuthorRepository interface {
	logWithCtx(ctx context.Context, function string) *log.Entry
	Create(ctx context.Context, author *model.Author) (*model.Author, error)
	GetByIDs(ctx context.Context, ids ...uint) (*[]model.Author, error)
	Update(ctx context.Context, author model.Author) error
	DeleteById(ctx context.Context, id uint) error
	GetAll(ctx context.Context) (*[]model.Author, error)
	GetAuthorsBook(ctx context.Context, author *model.Author) (*[]model.Book, error)
}
