package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Fine struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	LoanID    uuid.UUID
	Loan      Loan
	MemberID  uuid.UUID
	Amount    float64
	Reason    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
