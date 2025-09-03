package service

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/nanoLeinz/librarium/model/dto"
)

type AuthorService interface {
	logWithCtx(ctx context.Context, function string) *log.Entry
	Create(ctx context.Context, data *dto.AuthorRequest) (*dto.AuthorResponse, error)
	Update(ctx context.Context, id uint, data *dto.AuthorRequest) error
	DeleteById(ctx context.Context, id uint) error
	GetByIDs(ctx context.Context, ids ...uint) (*[]dto.AuthorResponse, error)
	GetAll(ctx context.Context) (*[]dto.AuthorResponse, error)
	GetAuthorsBook(ctx context.Context, author *dto.AuthorRequest) (*[]dto.BookResponseShort, error)
}
