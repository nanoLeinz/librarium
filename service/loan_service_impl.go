package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/internal/enum"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/repository"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LoanServiceImpl struct {
	log        *log.Logger
	loanRepo   repository.LoanRepository
	memberRepo repository.MemberRepository
	copyRepo   repository.BookCopyRepository
}

func NewLoanServiceImpl(log *log.Logger, loanRepo repository.LoanRepository, memberRepo repository.MemberRepository, copyRepo repository.BookCopyRepository) LoanService {
	return &LoanServiceImpl{
		log:        log,
		loanRepo:   loanRepo,
		memberRepo: memberRepo,
		copyRepo:   copyRepo,
	}
}

func (s *LoanServiceImpl) logWithCtx(ctx context.Context, function string) *log.Entry {

	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

	logger := s.log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": function,
	})

	return logger
}
func (s *LoanServiceImpl) Create(ctx context.Context, data *dto.LoanRequest) (*dto.LoanResponse, error) {
	logger := s.logWithCtx(ctx, "LoanService.Create").
		WithFields(log.Fields{
			"memberID":   data.MemberID,
			"bookCopyID": data.BookCopyID,
		})

	logger.Info("received create loan request")

	member, err := s.memberRepo.GetByID(ctx, data.MemberID)
	if err != nil {
		logger.WithError(err).Error("failed to get member by ID")
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, myerror.NewNotFoundError("member")
		default:
			return nil, myerror.InternalServerErr
		}
	}

	if member.AccountStatus == enum.SuspendedAccount.String() || member.AccountStatus != enum.ActiveAccount.String() {
		logger.Warn("member account is not active or is suspended")
		return nil, myerror.NewBadRequestError("account suspended")
	}

	result, err := s.copyRepo.GetByID(ctx, data.BookCopyID)
	if err != nil {
		logger.WithError(err).Error("failed to get book copy by ID")
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, myerror.NewNotFoundError("book copy")
		default:
			return nil, myerror.InternalServerErr
		}
	}

	if result.Status != enum.AvailableCopy.String() {
		logger.Warn("chosen book copy is unavailable")
		return nil, myerror.NewBadRequestError("chosen copy unavailable")
	}

	loan := model.Loan{
		MemberID:   data.MemberID,
		BookCopyID: data.BookCopyID,
		LoanDate:   time.Now(),
		DueDate:    time.Now().Add(time.Hour * 24 * 7),
		Status:     enum.ActiveLoan.String(),
	}

	query, err := s.loanRepo.Create(ctx, &loan)
	if err != nil {
		logger.WithError(err).Error("failed to create loan in repository")
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			logger.WithError(err).Error("duplicate loan error")
			return nil, myerror.NewDuplicateError("loan")
		}
		return nil, myerror.InternalServerErr
	}

	log.Info("updating copy status to loaned")
	result.Status = enum.LoanedCopy.String()
	_ = s.copyRepo.Update(ctx, result)

	logger.WithField("loanID", query.ID).Info("loan created successfully")
	response := dto.ToLoanResponse(*query)
	return &response, nil
}

func (s *LoanServiceImpl) Update(ctx context.Context, id uuid.UUID, data *dto.LoanUpdateRequest) error {
	logger := s.logWithCtx(ctx, "LoanService.Update").
		WithField("loanID", id)
	logger.Info("received update loan request")

	var status string
	if strings.ToLower(data.Status) == enum.ReturnedLoan.String() {
		status = enum.ReturnedLoan.String()
	} else if strings.ToLower(data.Status) == enum.OverdueLoan.String() {
		status = enum.OverdueLoan.String()
	} else {
		return myerror.NewBadRequestError("status invalid")
	}

	loan := model.Loan{
		ID:     id,
		Status: status,
	}

	err := s.loanRepo.Update(ctx, &loan)
	if err != nil {
		logger.WithError(err).Error("failed to update loan in repository")
		switch err {
		case gorm.ErrRecordNotFound:
			return myerror.NewNotFoundError("loan")
		default:
			return myerror.InternalServerErr
		}
	}

	log.Infof("updating copy status to %s", data.BookStatus)
	copy := model.BookCopy{
		Model:  gorm.Model{ID: data.BookCopyID},
		Status: data.BookStatus,
	}

	_ = s.copyRepo.Update(ctx, &copy)

	logger.Info("loan updated successfully")
	return nil
}

func (s *LoanServiceImpl) DeleteById(ctx context.Context, id uuid.UUID) error {
	logger := s.logWithCtx(ctx, "LoanService.DeleteById").
		WithField("loanID", id)

	logger.Info("received delete loan request")

	err := s.loanRepo.DeleteByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to delete loan in repository")
		switch err {
		case gorm.ErrRecordNotFound:
			return myerror.NewNotFoundError("loan")
		default:
			return myerror.InternalServerErr
		}
	}

	logger.Info("loan deleted successfully")
	return nil
}

func (s *LoanServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*dto.LoanResponse, error) {
	logger := s.logWithCtx(ctx, "LoanService.GetByID").
		WithField("loanID", id)

	logger.Info("received get loan by ID request")

	result, err := s.loanRepo.GetByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to get loan by ID from repository")
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, myerror.NewNotFoundError("loan")
		default:
			return nil, myerror.InternalServerErr
		}
	}
	response := dto.ToLoanResponse(*result)
	logger.Info("loan fetched successfully")
	return &response, nil
}

func (s *LoanServiceImpl) GetAll(ctx context.Context) (*[]dto.LoanResponse, error) {
	logger := s.logWithCtx(ctx, "LoanService.GetAll")
	logger.Info("received get all loans request")

	result, err := s.loanRepo.GetAll(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to get all loans from repository")
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, myerror.NewNotFoundError("loan")
		default:
			return nil, myerror.InternalServerErr
		}
	}

	responses := []dto.LoanResponse{}
	for _, v := range *result {
		responses = append(responses, dto.ToLoanResponse(v))
	}

	logger.WithField("count", len(responses)).Info("all loans fetched successfully")
	return &responses, nil
}
