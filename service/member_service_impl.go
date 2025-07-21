package service

import (
	"context"
	"errors"
	"fmt"
	"log" // Add logging

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
	log.Println("NewMemberServiceImpl: initializing service") // Logging
	return &MemberServiceImpl{repo: repo}
}

func (s MemberServiceImpl) GetAllMembers(ctx context.Context) ([]dto.MemberResponse, error) {
	log.Println("GetAllMembers: fetching all members") // Logging

	result, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("GetAllMembers: failed to fetch members: %v", err) // Logging
		return nil, err
	}
	log.Printf("GetAllMembers: fetched %d members", len(*result)) // Logging

	members := make([]dto.MemberResponse, 0, len(*result))
	for _, v := range *result {
		members = append(members, dto.ToMemberResponse(v))
	}
	log.Println("GetAllMembers: converted members to response DTOs") // Logging

	return members, nil
}

func (s MemberServiceImpl) CreateMember(ctx context.Context, data *dto.MemberCreateRequest) (*dto.MemberResponse, error) {
	log.Printf("CreateMember: received data: %+v", data) // Logging

	hashedpass, err := helper.HashPassword(data.Password)
	if err != nil {
		log.Printf("CreateMember: failed to hash password: %v", err) // Logging
		return nil, errors.New("failed hashing password")
	}
	log.Println("CreateMember: password hashed successfully") // Logging

	user := model.Member{
		Email:         data.Email,
		Password:      hashedpass,
		FullName:      data.FullName,
		AccountStatus: data.AccountStatus,
		Role:          data.Role,
	}
	log.Printf("CreateMember: member model prepared: %+v", user) // Logging

	result, err := s.repo.Create(ctx, &user)
	if err != nil {
		log.Printf("CreateMember: failed to create member in repository: %v", err) // Logging
		return nil, err
	}
	log.Printf("CreateMember: member created in repository: %+v", result) // Logging

	member := dto.ToMemberResponse(*result)
	log.Printf("CreateMember: converted to response DTO: %+v", member) // Logging

	return &member, nil
}

func (s MemberServiceImpl) UpdateMember(ctx context.Context, data *dto.MemberUpdateRequest) error {
	log.Printf("UpdateMember: received update request: %+v", data) // Logging

	updates := dto.StructToMap(data)
	log.Printf("UpdateMember: converted update request to map: %+v", updates) // Logging

	err := s.repo.Update(ctx, *data.ID, &updates)
	if err != nil {
		log.Printf("UpdateMember: failed to update member: %v", err) // Logging
		return fmt.Errorf("update data failed: %w", err)
	}

	log.Println("UpdateMember: member updated successfully") // Logging
	return nil
}

func (s MemberServiceImpl) GetMemberByID(ctx context.Context, id uuid.UUID) (*dto.MemberResponse, error) {
	log.Printf("GetMemberByID: fetching member with ID: %s", id) // Logging

	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("GetMemberByID: failed to get member: %v", err) // Logging
		return nil, fmt.Errorf("failed to get record %w", err)
	}

	result := dto.ToMemberResponse(*data)
	log.Printf("GetMemberByID: converted to response DTO: %+v", result) // Logging

	return &result, nil
}

func (s MemberServiceImpl) DeleteMemberByID(ctx context.Context, id uuid.UUID) error {
	log.Printf("DeleteMemberByID: deleting member with ID: %s", id) // Logging

	err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		log.Printf("DeleteMemberByID: failed to delete member: %v", err) // Logging
		return fmt.Errorf("failed to delete record %w", err)
	}

	log.Println("DeleteMemberByID: member deleted successfully") // Logging
	return nil
}
