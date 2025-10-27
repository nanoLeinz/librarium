package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model/dto"
	log "github.com/sirupsen/logrus"
)

type ReservationService interface {
	logWithCtx(ctx context.Context, function string) *log.Entry
	Create(ctx context.Context, data *dto.ReservationRequest) (*dto.ReservationResponse, error)
	Update(ctx context.Context, id uuid.UUID, data *dto.ReservationRequest) error
	DeleteById(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.ReservationResponse, error)
	GetAll(ctx context.Context) ([]dto.ReservationResponse, error)
}
