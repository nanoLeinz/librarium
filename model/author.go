package model

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	Name      string
	Biography string
	BirthYear time.Time
	Book      []Book `gorm:"many2many:author_books;"`
}
