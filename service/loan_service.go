package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model/dto"
	log "github.com/sirupsen/logrus"
)

type LoanService interface {
	logWithCtx(ctx context.Context, function string) *log.Entry
	Create(ctx context.Context, data *dto.LoanRequest) (*dto.LoanResponse, error)
	Update(ctx context.Context, id uuid.UUID, data *dto.LoanUpdateRequest) error
	DeleteById(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.LoanResponse, error)
	GetAll(ctx context.Context) (*[]dto.LoanResponse, error)
}
