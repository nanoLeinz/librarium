package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Member struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email         string
	Password      string
	FullName      string
	Role          string
	AccountStatus string
	Loan          []Loan
	Fine          []Fine
	Reservation   []Reservation
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
