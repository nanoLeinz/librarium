package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title           string
	ISBN            string
	PublicationYear int
	Genre           string
	Author          []Author `gorm:"many2many:author_books;"`
	BookCopy        []BookCopy
	Reservation     []Reservation
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}
