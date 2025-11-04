package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/internal/helper"
	"github.com/nanoLeinz/librarium/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LoanRepositoryImpl struct {
	log *log.Logger
	db  *gorm.DB
}

func NewLoanRepository(log *log.Logger, db *gorm.DB) LoanRepository {
	return &LoanRepositoryImpl{
		log: log,
		db:  db,
	}
}

func (s *LoanRepositoryImpl) logWithCtx(ctx context.Context, function string) *log.Entry {

	logger := s.log.WithFields(log.Fields{
		"traceID":  ctx.Value(helper.KeyCon("traceID")).(string),
		"function": function,
	})

	return logger
}

func (s *LoanRepositoryImpl) Create(ctx context.Context, loan *model.Loan) (*model.Loan, error) {

	logger := s.logWithCtx(ctx, "LoanRepository.Create").
		WithFields(log.Fields{
			"memberID":   loan.MemberID,
			"bookCopyID": loan.BookCopyID,
		})

	logger.Info("executing query")

	result := s.db.WithContext(ctx).Create(loan)

	if result.Error != nil {
		logger.WithError(result.Error).Error("failed executing query")
		return nil, result.Error
	}

	logger.WithField("loanID", loan.ID).Info("successfully executed query")

	return loan, nil
}
func (s *LoanRepositoryImpl) Update(ctx context.Context, loan *model.Loan) error {

	logger := s.logWithCtx(ctx, "LoanRepository.Update").
		WithFields(log.Fields{
			"loanID": loan.ID,
			"status": loan.Status,
		})

	logger.Info("executing query")

	result := s.db.WithContext(ctx).Updates(loan)

	if err := result.Error; err != nil {
		logger.WithError(err).Error("failed executing query")
		return err
	}

	logger.Info("query executed successfully")

	return nil

}
func (s *LoanRepositoryImpl) DeleteByID(ctx context.Context, loanID uuid.UUID) error {

	logger := s.logWithCtx(ctx, "LoanRepository.DeleteByID").
		WithFields(log.Fields{
			"loanID": loanID,
		})

	logger.Info("executing query")

	result := s.db.WithContext(ctx).Delete(&model.Loan{}, loanID)

	if result.Error != nil {
		logger.WithError(result.Error).Error("failed executing query")
		return result.Error
	} else if result.RowsAffected == 0 {
		logger.Debug("query executed but 0 rows affected")
	} else {
		logger.Info("query executed successfully")
	}

	return nil
}
func (s *LoanRepositoryImpl) GetByID(ctx context.Context, loanID uuid.UUID) (*model.Loan, error) {

	logger := s.logWithCtx(ctx, "LoanRepository.loanIDs").
		WithFields(log.Fields{
			"loanIDs": loanID,
		})

	logger.Info("executing query")

	var loans = model.Loan{}
	q := s.db.WithContext(ctx).Find(&loans, loanID)
	if err := q.Error; err != nil {
		logger.WithError(err).Error("failed executing query")
		return nil, err
	} else if q.RowsAffected == 0 {
		logger.WithError(gorm.ErrRecordNotFound).Error("record not found")
		return nil, gorm.ErrRecordNotFound
	}

	logger.Info("query executed successfully")
	return &loans, nil

}
func (s *LoanRepositoryImpl) GetAll(ctx context.Context) (*[]model.Loan, error) {
	s.logWithCtx(ctx, "LoanRepository.GetAll").Info("executing query")

	var loans = []model.Loan{}

	if err := s.db.WithContext(ctx).Scopes(helper.Paginator(ctx)).Find(&loans).Error; err != nil {
		s.logWithCtx(ctx, "LoanRepository.GetAll").
			WithError(err).
			Error("failed executing query")

		return nil, err
	} else if len(loans) == 0 {
		s.logWithCtx(ctx, "LoanRepository.GetAll").
			Debug("query executed but data not found")

		return nil, gorm.ErrRecordNotFound
	} else {
		s.logWithCtx(ctx, "LoanRepository.GetAll").
			Info("query executed successfully")
	}

	return &loans, nil

}
