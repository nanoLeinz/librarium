package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MemberServiceImpl struct {
	repo repository.MemberRepository
	log  *logrus.Logger
}

func NewMemberServiceImpl(repo repository.MemberRepository, log *logrus.Logger) MemberService {
	return &MemberServiceImpl{
		repo: repo,
		log:  log,
	}
}

func (s MemberServiceImpl) GetAllMembers(ctx context.Context) ([]dto.MemberResponse, error) {
	s.log.WithField("function", "GetAllMembers").Info("Attempting to fetch all members")

	result, err := s.repo.GetAll(ctx)
	if err != nil {
		s.log.WithField("function", "GetAllMembers").WithError(err).Error("Failed to fetch members from repository")
		return nil, err
	}

	members := make([]dto.MemberResponse, 0, len(*result))
	for _, v := range *result {
		members = append(members, dto.ToMemberResponse(v))
	}

	s.log.WithFields(logrus.Fields{
		"function": "GetAllMembers",
		"count":    len(members),
	}).Info("Successfully fetched and converted all members")

	return members, nil
}

func (s MemberServiceImpl) CreateMember(ctx context.Context, data *dto.MemberCreateRequest) (*dto.MemberResponse, error) {
	s.log.WithFields(logrus.Fields{
		"function": "CreateMember",
		"email":    data.Email,
	}).Info("Attempting to create a new member")

	hashedpass, err := helper.HashPassword(data.Password)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"function": "CreateMember",
			"email":    data.Email,
		}).WithError(err).Error("Failed to hash password")
		return nil, errors.New("failed hashing password")
	}

	user := model.Member{
		Email:         data.Email,
		Password:      hashedpass,
		FullName:      data.FullName,
		AccountStatus: data.AccountStatus,
		Role:          data.Role,
	}

	result, err := s.repo.Create(ctx, &user)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"function": "CreateMember",
			"email":    data.Email,
		}).WithError(err).Error("Failed to create member in repository")

		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				s.log.WithFields(logrus.Fields{
					"function": "CreateMember",
					"email":    data.Email,
				}).WithError(err).Error("duplicated constraints for member")
				return nil, myerror.NewDuplicateError("member")
			}

			s.log.WithFields(logrus.Fields{
				"function": "CreateMember",
				"email":    data.Email,
			}).WithError(err).Error("failed to craete member")
			return nil, myerror.InternalServerErr
		}
	}

	member := dto.ToMemberResponse(*result)
	s.log.WithFields(logrus.Fields{
		"function": "CreateMember",
		"memberID": member.ID,
	}).Info("Successfully created new member")

	return &member, nil
}

func (s MemberServiceImpl) UpdateMember(ctx context.Context, data *dto.MemberUpdateRequest) error {
	s.log.WithFields(logrus.Fields{
		"function": "UpdateMember",
		"memberID": data.ID,
	}).Info("Attempting to update member")

	updates := dto.StructToMap(*data)

	err := s.repo.Update(ctx, data.ID, &updates)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"function": "UpdateMember",
			"memberID": data.ID,
		}).WithError(err).Error("Failed to update member in repository")

		switch err {
		case gorm.ErrRecordNotFound:
			s.log.WithError(err).Error("member not found")
			return myerror.NewNotFoundError("member")
		default:
			s.log.WithError(err).Error(err.Error())
			return myerror.InternalServerErr
		}

	}

	s.log.WithFields(logrus.Fields{
		"function": "UpdateMember",
		"memberID": data.ID,
	}).Info("Successfully updated member")
	return nil
}

func (s MemberServiceImpl) GetMemberByID(ctx context.Context, id uuid.UUID) (*dto.MemberResponse, error) {
	s.log.WithFields(logrus.Fields{
		"function": "GetMemberByID",
		"memberID": id,
	}).Info("Attempting to fetch member by ID")

	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"function": "GetMemberByID",
			"memberID": id,
		}).WithError(err).Error("Failed to get member from repository")

		switch err {
		case gorm.ErrRecordNotFound:
			return nil, myerror.NewNotFoundError("member")
		default:
			return nil, myerror.InternalServerErr
		}

	}

	result := dto.ToMemberResponse(*data)
	s.log.WithFields(logrus.Fields{
		"function": "GetMemberByID",
		"memberID": id,
	}).Info("Successfully fetched member by ID")

	return &result, nil
}

func (s MemberServiceImpl) GetMemberByEmail(ctx context.Context, email string) (*dto.MemberResponse, error) {
	s.log.WithFields(logrus.Fields{
		"function": "GetMemberByEmail",
		"email":    email,
	}).Info("Attempting to fetch member by email")

	data, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"function": "GetMemberByEmail",
			"email":    email,
		}).WithError(err).Error("Failed to get member from repository")

		switch err {
		case gorm.ErrRecordNotFound:
			return nil, myerror.NewNotFoundError("member")
		default:
			return nil, myerror.InternalServerErr
		}

	}

	result := dto.ToMemberResponse(*data)
	s.log.WithFields(logrus.Fields{
		"function": "GetMemberByEmail",
		"email":    email,
		"memberID": result.ID,
	}).Info("Successfully fetched member by email")

	return &result, nil
}

func (s MemberServiceImpl) DeleteMemberByID(ctx context.Context, id uuid.UUID) error {
	s.log.WithFields(logrus.Fields{
		"function": "DeleteMemberByID",
		"memberID": id,
	}).Info("Attempting to delete member by ID")

	err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"function": "DeleteMemberByID",
			"memberID": id,
		}).WithError(err).Error("Failed to delete member in repository")
		switch err {
		case gorm.ErrRecordNotFound:
			return myerror.NewNotFoundError("member")
		default:
			return myerror.InternalServerErr
		}
	}

	s.log.WithFields(logrus.Fields{
		"function": "DeleteMemberByID",
		"memberID": id,
	}).Info("Successfully deleted member")
	return nil
}
