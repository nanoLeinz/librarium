package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/internal/enum"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/repository"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ReservationServiceImpl struct {
	repo       repository.ReservationRepository
	memberRepo repository.MemberRepository
	log        *log.Logger
}

func NewReservationService(repo repository.ReservationRepository, log *log.Logger) ReservationService {
	return &ReservationServiceImpl{
		repo: repo,
		log:  log,
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

	member, err := s.memberRepo.GetByID(ctx, data.MemberID)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, myerror.NewNotFoundError("member")
		default:
			return nil, myerror.InternalServerErr
		}
	}
	if member.AccountStatus == enum.SuspendedAccount.String() || member.AccountStatus != enum.ActiveAccount.String() {
		log.Warn("member account is not active or is suspended")
		return nil, myerror.NewBadRequestError("account suspended")
	}

	queue := s.repo.GetLastQueue(ctx, data.BookID)

	reservation := &model.Reservation{
		BookID:          data.BookID,
		MemberID:        data.MemberID,
		Status:          enum.PendingReserv.String(),
		ReservationDate: time.Now().Local(),
		QueuePosition:   queue,
	}

	result, err := s.repo.Create(ctx, reservation)
	if err != nil {
		return nil, myerror.InternalServerErr
	}

	res := dto.ToReservationResponse(*result)

	return &res, nil
}

func (s *ReservationServiceImpl) Update(ctx context.Context, id uuid.UUID, data *dto.ReservationRequest) error {
	return nil
}
func (s *ReservationServiceImpl) DeleteById(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (s *ReservationServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*[]dto.ReservationResponse, error) {
	return nil, nil
}
func (s *ReservationServiceImpl) GetAll(ctx context.Context) (*[]dto.ReservationResponse, error) {
	return nil, nil
}
