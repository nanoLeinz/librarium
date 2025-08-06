package dto

import (
	"github.com/google/uuid"
)

type BookCreateReq struct {
	Title           string     `json:"title"`
	ISBN            string     `json:"isbn"`
	PublicationYear int        `json:"publication_year"`
	Genre           string     `json:"genre"`
	AuthorIds       uuid.UUIDs `json:"authors"`
}
