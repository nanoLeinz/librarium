package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
)

type MemberResponse struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	FullName      string    `json:"full_name"`
	Role          string    `json:"-"`
	AccountStatus string    `json:"account_status"`
	CreatedAt     time.Time `json:"created_at"`
}

func ToMemberResponse(member model.Member) MemberResponse {
	return MemberResponse{
		ID:            member.ID,
		Email:         member.Email,
		FullName:      member.FullName,
		Role:          member.Role,
		AccountStatus: member.AccountStatus,
		CreatedAt:     member.CreatedAt,
		Password:      member.Password,
	}

}
