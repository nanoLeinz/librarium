package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
	log "github.com/sirupsen/logrus"
)

type LoanRepository interface {
	logWithCtx(ctx context.Context, function string) *log.Entry
	Create(ctx context.Context, loan *model.Loan) (*model.Loan, error)
	Update(ctx context.Context, loan *model.Loan) error
	DeleteByID(ctx context.Context, loanID uuid.UUID) error
	GetByID(ctx context.Context, loanIDs uuid.UUID) (*model.Loan, error)
	GetAll(ctx context.Context) (*[]model.Loan, error)
}
