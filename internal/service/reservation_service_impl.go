package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/internal/enum"
	"github.com/nanoLeinz/librarium/internal/helper"
	"github.com/nanoLeinz/librarium/internal/model"
	"github.com/nanoLeinz/librarium/internal/model/dto"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/internal/repository"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ReservationServiceImpl struct {
	repo       repository.ReservationRepository
	memberRepo repository.MemberRepository
	log        *log.Logger
}

func NewReservationService(log *log.Logger, repo repository.ReservationRepository, memberRepo repository.MemberRepository) ReservationService {
	return &ReservationServiceImpl{
		repo:       repo,
		log:        log,
		memberRepo: memberRepo,
	}
}

func (s *ReservationServiceImpl) logWithCtx(ctx context.Context, function string) *log.Entry {

	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

	logger := s.log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": function,
	})

	return logger
}

func (s *ReservationServiceImpl) Create(ctx context.Context, data *dto.ReservationRequest) (*dto.ReservationResponse, error) {
	logger := s.logWithCtx(ctx, "ReservationService.Create").
		WithFields(log.Fields{
			"bookID":   data.BookID,
			"memberID": data.MemberID,
		})

	logger.Info("received create reservation request")

	member, err := s.memberRepo.GetByID(ctx, data.MemberID)
	if err != nil {
		logger.WithError(err).Error("failed to fetch member")
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, myerror.NewNotFoundError("member")
		default:
			return nil, myerror.InternalServerErr
		}
	}

	if member.AccountStatus == enum.SuspendedAccount.String() || member.AccountStatus != enum.ActiveAccount.String() {
		logger.WithField("accountStatus", member.AccountStatus).Warn("member account is not active or is suspended")
		return nil, myerror.NewBadRequestError("account suspended")
	}

	queue := s.repo.GetLastQueue(ctx, data.BookID)

	reservation := &model.Reservation{
		BookID:          data.BookID,
		MemberID:        data.MemberID,
		Status:          enum.PendingReserv.String(),
		ReservationDate: time.Now().Local(),
		QueuePosition:   queue + 1,
	}

	result, err := s.repo.Create(ctx, reservation)
	if err != nil {
		logger.WithError(err).Error("failed to create reservation in repository")
		return nil, myerror.InternalServerErr
	}

	logger.WithFields(log.Fields{
		"reservationID": result.ID,
		"queuePosition": result.QueuePosition,
	}).Info("reservation created successfully")

	res := dto.ToReservationResponse(*result)
	return &res, nil
}

func (s *ReservationServiceImpl) Update(ctx context.Context, id uuid.UUID, data *dto.ReservationRequest) error {
	logger := s.logWithCtx(ctx, "ReservationService.Update").
		WithFields(log.Fields{
			"reservationID": id,
			"bookID":        data.BookID,
			"memberID":      data.MemberID,
			"status":        data.Status,
		})

	logger.Info("received update reservation request")

	req := model.Reservation{
		ID:        id,
		Status:    data.Status,
		UpdatedAt: time.Now().Local(),
	}

	err := s.repo.Update(ctx, req)
	if err != nil {
		logger.WithError(err).Error("failed to update reservation in repository")
		switch err {
		case gorm.ErrRecordNotFound:
			return myerror.NewNotFoundError("reservation")
		default:
			return myerror.InternalServerErr
		}
	}

	logger.Info("reservation updated successfully")

	if req.Status != "" {
		logger.WithFields(log.Fields{
			"reservationID": req.ID,
			"bookID":        req.BookID,
		}).Info("updating related queue due to status change")
		_ = s.repo.UpdateRelatedQueue(ctx, id)
	}

	return nil
}

func (s *ReservationServiceImpl) DeleteById(ctx context.Context, id uuid.UUID) error {
	logger := s.logWithCtx(ctx, "ReservationService.DeleteById").
		WithField("reservationID", id)

	logger.Info("received delete reservation request")

	err := s.repo.DeleteById(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to delete reservation in repository")
		switch err {
		case gorm.ErrRecordNotFound:
			return myerror.NewNotFoundError("reservation")
		default:
			return myerror.InternalServerErr
		}
	}

	logger.Info("reservation deleted successfully")
	return nil
}

func (s *ReservationServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*dto.ReservationResponse, error) {
	logger := s.logWithCtx(ctx, "ReservationService.GetByID").
		WithField("reservationID", id)

	logger.Info("received get reservation by ID request")

	res, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to fetch reservation from repository")
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, myerror.NewNotFoundError("reservation")
		default:
			return nil, myerror.InternalServerErr
		}
	}

	response := dto.ToReservationResponse(*res)
	logger.WithField("reservationID", response.ID).Info("reservation fetched successfully")
	return &response, nil
}

func (s *ReservationServiceImpl) GetAll(ctx context.Context) ([]dto.ReservationResponse, error) {
	logger := s.logWithCtx(ctx, "ReservationService.GetAll")
	logger.Info("received get all reservations request")

	res, err := s.repo.GetAll(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to fetch reservations from repository")
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, myerror.NewNotFoundError("reservation")
		default:
			return nil, myerror.InternalServerErr
		}
	}

	response := []dto.ReservationResponse{}
	for _, v := range res {
		response = append(response, dto.ToReservationResponse(v))
	}

	logger.WithField("count", len(response)).Info("all reservations fetched successfully")
	return response, nil
}
