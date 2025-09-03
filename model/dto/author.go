package dto

import (
	"time"

	"github.com/nanoLeinz/librarium/model"
)

type AuthorRequest struct {
	ID        uint   `json:"-"`
	Name      string `json:"name"`
	Biography string `json:"biography"`
	BirthYear int    `json:"year_of_birth"`
}

type AuthorResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Biography string    `json:"biography"`
	BirthYear int       `json:"year_of_birth"`
	CreatedAt time.Time `json:"created_at"`
}

// type AuthorShort struct {
// 	ID   uuid.UUID `json:"id"`
// 	Name string    `json:"name"`
// }

func ToAuthorResponse(data model.Author) AuthorResponse {
	return AuthorResponse{
		ID:        data.ID,
		Name:      data.Name,
		Biography: data.Biography,
		BirthYear: data.BirthYear,
		CreatedAt: data.CreatedAt,
	}
}
