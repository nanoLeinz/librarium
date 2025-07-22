package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model/dto"
)

type MemberService interface {
	GetAllMembers(ctx context.Context) ([]dto.MemberResponse, error)
	CreateMember(ctx context.Context, data *dto.MemberCreateRequest) (*dto.MemberResponse, error)
	UpdateMember(ctx context.Context, data *dto.MemberUpdateRequest) error
	GetMemberByID(ctx context.Context, id uuid.UUID) (*dto.MemberResponse, error)
	GetMemberByEmail(ctx context.Context, email string) (*dto.MemberResponse, error)
	DeleteMemberByID(ctx context.Context, id uuid.UUID) error
}
