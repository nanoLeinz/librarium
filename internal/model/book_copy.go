package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookCopy struct {
	gorm.Model
	Status string
	BookID uuid.UUID
	Loan   []Loan
}
