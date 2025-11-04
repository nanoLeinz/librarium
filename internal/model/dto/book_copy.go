package dto

import "github.com/google/uuid"

type BookCopyRequest struct {
	Copies uint      `json:"copies"`
	Status string    `json:"status"`
	BookID uuid.UUID `json:"book_id"`
}

type BookCopyResponse struct {
	ID     uint      `json:"id"`
	Status string    `json:"status"`
	BookID uuid.UUID `json:"book_id"`
}
