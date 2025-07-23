package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MemberRepositoryImpl struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewMemberRepository(db *gorm.DB, log *logrus.Logger) MemberRepository {
	return &MemberRepositoryImpl{
		db:  db,
		log: log,
	}
}

func (s *MemberRepositoryImpl) Create(ctx context.Context, data *model.Member) (*model.Member, error) {
	s.log.WithFields(logrus.Fields{
		"function": "Create",
		"memberID": data.ID,
	}).Info("Attempting to create a new member")

	result := s.db.WithContext(ctx).Create(data)
	if result.Error != nil {
		s.log.WithFields(logrus.Fields{
			"function": "Create",
			"memberID": data.ID,
		}).WithError(result.Error).Error("Failed to create member")
		return nil, result.Error
	}

	s.log.WithFields(logrus.Fields{
		"function": "Create",
		"memberID": data.ID,
	}).Info("Member created successfully")
	return data, nil
}

func (s *MemberRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	s.log.WithFields(logrus.Fields{
		"function": "DeleteByID",
		"memberID": id.String(),
	}).Info("Attempting to delete member")

	result := s.db.WithContext(ctx).Delete(&model.Member{}, id)
	if result.Error != nil {
		s.log.WithFields(logrus.Fields{
			"function": "DeleteByID",
			"memberID": id.String(),
		}).WithError(result.Error).Error("Failed to delete member")
		return result.Error
	}

	if result.RowsAffected == 0 {
		s.log.WithFields(logrus.Fields{
			"function": "DeleteByID",
			"memberID": id.String(),
		}).Warn("Attempted to delete a member that does not exist")
	} else {
		s.log.WithFields(logrus.Fields{
			"function": "DeleteByID",
			"memberID": id.String(),
		}).Info("Member deleted successfully")
	}

	return nil
}

func (s *MemberRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Member, error) {
	s.log.WithFields(logrus.Fields{
		"function": "GetByID",
		"memberID": id.String(),
	}).Info("Attempting to fetch member by ID")

	var data model.Member
	result := s.db.WithContext(ctx).First(&data, id)

	if result.Error != nil {
		s.log.WithFields(logrus.Fields{
			"function": "GetByID",
			"memberID": id.String(),
		}).WithError(result.Error).Error("Failed to fetch member by ID")
		return nil, result.Error
	}

	s.log.WithFields(logrus.Fields{
		"function": "GetByID",
		"memberID": id.String(),
	}).Info("Member fetched successfully")
	return &data, nil
}

func (s *MemberRepositoryImpl) GetByEmail(ctx context.Context, email string) (*model.Member, error) {
	s.log.WithFields(logrus.Fields{
		"function": "GetByEmail",
		"email":    email,
	}).Info("Attempting to fetch member by email")

	data := &model.Member{}
	result := s.db.WithContext(ctx).Where("email = ?", email).First(data)

	if result.Error != nil {
		s.log.WithFields(logrus.Fields{
			"function": "GetByEmail",
			"email":    email,
		}).WithError(result.Error).Error("Failed to fetch member by email")
		return nil, result.Error
	}

	s.log.WithFields(logrus.Fields{
		"function": "GetByEmail",
		"email":    email,
		"memberID": data.ID,
	}).Info("Member fetched successfully")
	return data, nil
}

func (s *MemberRepositoryImpl) GetAll(ctx context.Context) (*[]model.Member, error) {
	s.log.WithField("function", "GetAll").Info("Attempting to fetch all members")

	var data []model.Member
	result := s.db.WithContext(ctx).Find(&data)

	if result.Error != nil {
		s.log.WithField("function", "GetAll").WithError(result.Error).Error("Failed to fetch all members")
		return nil, result.Error
	}

	s.log.WithFields(logrus.Fields{
		"function": "GetAll",
		"count":    len(data),
	}).Info("All members fetched successfully")
	return &data, nil
}

func (s *MemberRepositoryImpl) Update(ctx context.Context, id uuid.UUID, data *map[string]interface{}) error {
	s.log.WithFields(logrus.Fields{
		"function": "Update",
		"memberID": id.String(),
		"data":     *data,
	}).Info("Attempting to update member")

	result := s.db.WithContext(ctx).Model(&model.Member{}).Where("id = ?", id).Updates(data)

	if result.Error != nil {
		s.log.WithFields(logrus.Fields{
			"function": "Update",
			"memberID": id.String(),
		}).WithError(result.Error).Error("Failed to update member")
		return result.Error
	}

	if result.RowsAffected == 0 {
		s.log.WithFields(logrus.Fields{
			"function": "Update",
			"memberID": id.String(),
		}).Warn("Update operation did not affect any rows, member may not exist")
		return gorm.ErrRecordNotFound
	}

	s.log.WithFields(logrus.Fields{
		"function": "Update",
		"memberID": id.String(),
	}).Info("Member updated successfully")
	return nil
}
