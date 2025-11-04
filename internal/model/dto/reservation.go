package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/internal/model"
)

type ReservationRequest struct {
	BookID          uuid.UUID `json:"book_id"`
	MemberID        uuid.UUID `json:"member_id"`
	ReservationDate time.Time `json:"reservation_date"`
	Status          string    `json:"status"`
	QueuePosition   int       `json:"queue"`
}

type ReservationResponse struct {
	ID              uuid.UUID `json:"reservation_id"`
	BookID          uuid.UUID `json:"book_id"`
	ReservationDate time.Time `json:"reservation_date"`
	Status          string    `json:"status"`
	QueuePosition   int       `json:"queue"`
	CreatedAt       time.Time `json:"created_at"`
}

func ToReservationResponse(s model.Reservation) ReservationResponse {
	return ReservationResponse{
		ID:              s.ID,
		BookID:          s.BookID,
		ReservationDate: s.ReservationDate,
		Status:          s.Status,
		QueuePosition:   s.QueuePosition,
		CreatedAt:       s.CreatedAt,
	}
}
