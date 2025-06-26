package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookCopy struct {
	gorm.Model
	Status    string
	BookID    uuid.UUID
	Loan      []Loan
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
