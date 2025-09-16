package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
)

type LoanRequest struct {
	MemberID   uuid.UUID `json:"member_id"`
	BookCopyID uint      `json:"book_copy_id"`
	Status     string    `json:"status"`
}

type LoanUpdateRequest struct {
	Status     string `json:"status"`
	BookCopyID uint   `json:"book_copy_id"`
	BookStatus string `json:"book_status"`
}

type LoanResponse struct {
	ID         uuid.UUID `json:"id"`
	MemberID   uuid.UUID `json:"member_id"`
	BookCopyID uint      `json:"book_copy_id"`
	LoanDate   time.Time `json:"loan_date"`
	DueDate    time.Time `json:"due_date"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToLoanResponse(loan model.Loan) LoanResponse {
	return LoanResponse{
		ID:         loan.ID,
		MemberID:   loan.MemberID,
		BookCopyID: loan.BookCopyID,
		LoanDate:   loan.LoanDate,
		DueDate:    loan.DueDate,
		Status:     loan.Status,
		CreatedAt:  loan.CreatedAt,
	}
}
