package model

import (
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	Name      string `gorm:"uniqueIndex"`
	Biography string
	BirthYear int
	Book      []Book `gorm:"many2many:author_books;"`
}
