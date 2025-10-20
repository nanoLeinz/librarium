package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
	log "github.com/sirupsen/logrus"
)

type ReservationRepository interface {
	logWithCtx(ctx context.Context, function string) *log.Entry
	Create(ctx context.Context, reservation *model.Reservation) (*model.Reservation, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Reservation, error)
	Update(ctx context.Context, reservation model.Reservation) error
	DeleteById(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context) ([]model.Reservation, error)
	GetLastQueue(ctx context.Context, bookID uuid.UUID) int
	UpdateRelatedQueue(ctx context.Context, bookID uuid.UUID) error
}
