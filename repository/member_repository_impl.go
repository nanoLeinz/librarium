package repository

import (
	"context"
	"errors"
	"log" // Add logging

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
	"gorm.io/gorm"
)

type MemberRepositoryImpl struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	log.Println("NewMemberRepository: initializing repository") // Logging
	return &MemberRepositoryImpl{db: db}
}

func (s *MemberRepositoryImpl) Create(ctx context.Context, data *model.Member) (*model.Member, error) {
	log.Printf("Create: creating member with data: %+v", data) // Logging

	result := s.db.WithContext(ctx).Create(data)

	if result.Error != nil {
		log.Printf("Create: failed to create member: %v", result.Error) // Logging
		return nil, result.Error
	}

	log.Printf("Create: member created successfully: %+v", data) // Logging
	return data, nil
}

func (s *MemberRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	log.Printf("DeleteByID: deleting member with ID: %s", id) // Logging

	result := s.db.WithContext(ctx).Delete(&model.Member{}, id)

	if result.Error != nil {
		log.Printf("DeleteByID: failed to delete member: %v", result.Error) // Logging
	} else {
		log.Printf("DeleteByID: member deleted successfully: %s", id) // Logging
	}

	return result.Error
}

func (s *MemberRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Member, error) {
	log.Printf("GetByID: fetching member with ID: %s", id) // Logging

	var data *model.Member

	result := s.db.WithContext(ctx).First(data, id)

	if result.Error != nil {
		log.Printf("GetByID: failed to fetch member: %v", result.Error) // Logging
		return nil, result.Error
	}

	log.Printf("GetByID: member fetched successfully: %+v", data) // Logging
	return data, nil
}

func (s *MemberRepositoryImpl) GetByEmail(ctx context.Context, email string) (*model.Member, error) {
	log.Printf("GetByEmail: fetching member with email: %s", email) // Logging

	data := &model.Member{}

	result := s.db.WithContext(ctx).Where("email = ?", email).First(data)

	if result.Error != nil {
		log.Printf("GetByEmail: failed to fetch member: %v", result.Error) // Logging
		return nil, result.Error
	}

	log.Printf("GetByEmail: member fetched successfully: %+v", data) // Logging
	return data, nil
}

func (s *MemberRepositoryImpl) GetAll(ctx context.Context) (*[]model.Member, error) {
	log.Println("GetAll: fetching all members") // Logging

	data := []model.Member{}

	result := s.db.WithContext(ctx).Find(&data)

	if result.Error != nil {
		log.Printf("GetAll: failed to fetch members: %v", result.Error) // Logging
		return nil, result.Error
	}

	log.Printf("GetAll: members fetched successfully, count: %d", len(data)) // Logging
	return &data, nil
}

func (s *MemberRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data *map[string]interface{}) error {
	log.Printf("Update: updating member with ID: %s, data: %+v", id, data) // Logging

	result := s.db.WithContext(ctx).Model(&model.Member{}).Where("id = ?", id).Updates(data)

	if result.Error == gorm.ErrRecordNotFound {
		log.Printf("Update: member not found: %s", id) // Logging
		return errors.New("record tidak ditemukan")
	} else if result.Error != nil {
		log.Printf("Update: failed to update member: %v", result.Error) // Logging
		return result.Error
	}

	log.Printf("Update: member updated successfully: %s", id) // Logging
	return nil
}
