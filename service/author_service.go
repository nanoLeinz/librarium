package service

import (
	"context"

	"github.com/nanoLeinz/librarium/model/dto"
)

type AuthorService interface {
	Create(ctx context.Context, data *dto.AuthorRequest) (*dto.AuthorResponse, error)
}
