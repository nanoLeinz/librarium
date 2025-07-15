package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
	"gorm.io/gorm"
)

type MemberRepositoryImpl struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &MemberRepositoryImpl{db: db}
}

func (s *MemberRepositoryImpl) Create(ctx context.Context, data *model.Member) (*model.Member, error) {

	result := s.db.WithContext(ctx).Create(data)

	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func (s *MemberRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {

	result := s.db.WithContext(ctx).Delete(&model.Member{}, id)

	return result.Error
}

func (s *MemberRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Member, error) {

	var data *model.Member

	result := s.db.WithContext(ctx).First(data, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func (s *MemberRepositoryImpl) GetAll(ctx context.Context) (*[]model.Member, error) {

	var data []model.Member

	result := s.db.WithContext(ctx).Find(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	return &data, nil
}

func (s *MemberRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data *map[string]interface{}) error {

	result := s.db.WithContext(ctx).Model(&model.Member{}).Where("id = ?", id).Updates(data)

	if result.Error == gorm.ErrRecordNotFound {
		return errors.New("record tidak ditemukan")
	} else if result.Error != nil {
		return result.Error
	}

	return nil
}
