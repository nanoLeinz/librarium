package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
)

type MemberResponse struct {
	ID            uuid.UUID
	Email         string
	FullName      string
	AccountStatus string
	CreatedAt     time.Time
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
