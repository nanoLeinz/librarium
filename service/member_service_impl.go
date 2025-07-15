package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/model"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/repository"
)

type MemberServiceImpl struct {
	repo repository.MemberRepository
}

func NewMemberServiceImpl(repo repository.MemberRepository) MemberService {
	return &MemberServiceImpl{repo: repo}
}

func (s MemberServiceImpl) GetAllMembers(ctx context.Context) ([]dto.MemberResponse, error) {

	result, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}
	members := make([]dto.MemberResponse, 0, len(*result))

	for _, v := range *result {
		members = append(members, dto.ToMemberResponse(v))
	}

	return members, nil
}

func (s MemberServiceImpl) CreateMember(ctx context.Context, data *dto.MemberCreateRequest) (*dto.MemberResponse, error) {

	hashedpass, err := helper.HashPassword(data.Password)

	if err != nil {
		return nil, errors.New("failed hashing password")
	}

	user := model.Member{
		Email:         data.Email,
		Password:      hashedpass,
		FullName:      data.FullName,
		AccountStatus: data.AccountStatus,
	}

	result, err := s.repo.Create(ctx, &user)

	if err != nil {
		return nil, err
	}

	member := dto.ToMemberResponse(*result)

	return &member, nil
}

func (s MemberServiceImpl) UpdateMember(ctx context.Context, data *dto.MemberUpdateRequest) error {

	updates := dto.StructToMap(data)

	err := s.repo.Update(ctx, *data.ID, &updates)

	if err != nil {
		return fmt.Errorf("update data failed: %w", err)
	}

	return nil
}

func (s MemberServiceImpl) GetMemberByID(ctx context.Context, id uuid.UUID) (*dto.MemberResponse, error) {

	data, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("failed to get record %w", err)
	}

	result := dto.ToMemberResponse(*data)

	return &result, nil

}

func (s MemberServiceImpl) DeleteMemberByID(ctx context.Context, id uuid.UUID) error {

	err := s.repo.DeleteByID(ctx, id)

	if err != nil {
		return fmt.Errorf("failed to delete record %w", err)
	}

	return nil

}
