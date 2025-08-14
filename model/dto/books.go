package dto

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
)

type BookRequest struct {
	Title           string `json:"title"`
	ISBN            string `json:"isbn"`
	PublicationYear int    `json:"publication_year"`
	Genre           string `json:"genre"`
	InitialCopy     uint   `json:"initial_copy"`
	AuthorIds       []int  `json:"authors"`
}

type BookResponse struct {
	ID              uuid.UUID           `json:"id"`
	Title           string              `json:"title"`
	ISBN            string              `json:"isbn"`
	PublicationYear int                 `json:"publication_year"`
	Genre           string              `json:"genre"`
	Authors         []map[string]string `json:"authors"`
}

func ToBookResponse(book model.Book, authors []model.Author) BookResponse {

	var authorsSlice = []map[string]string{}

	for _, v := range authors {
		data := map[string]string{
			"author_id":   strconv.FormatUint(uint64(v.ID), 10),
			"author_name": v.Name,
		}

		authorsSlice = append(authorsSlice, data)
	}

	return BookResponse{
		ID:              book.ID,
		Title:           book.Title,
		ISBN:            book.ISBN,
		PublicationYear: book.PublicationYear,
		Genre:           book.Genre,
		Authors:         authorsSlice,
	}
}
