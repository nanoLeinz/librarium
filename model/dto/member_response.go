package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
)

type MemberResponse struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	AccountStatus string    `json:"account_status"`
	CreatedAt     time.Time `json:"created_at"`
}

func ToMemberResponse(member model.Member) MemberResponse {
	return MemberResponse{
		ID:            member.ID,
		Email:         member.Email,
		FullName:      member.FullName,
		AccountStatus: member.AccountStatus,
		CreatedAt:     member.CreatedAt,
	}

}
