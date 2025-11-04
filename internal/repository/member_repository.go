package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/internal/model"
)

type MemberRepository interface {
	Create(ctx context.Context, data *model.Member) (*model.Member, error)
	Update(ctx context.Context, id uuid.UUID, data *map[string]interface{}) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Member, error)
	GetByEmail(ctx context.Context, email string) (*model.Member, error)
	GetAll(ctx context.Context) (*[]model.Member, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
