package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Loan struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	MemberID   uuid.UUID
	BookCopyID uint
	LoanDate   time.Time
	DueDate    time.Time
	ReturnDate *time.Time
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
