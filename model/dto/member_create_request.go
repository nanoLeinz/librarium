package dto

type MemberCreateRequest struct {
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required"`
	FullName      string `json:"fullname" validate:"required"`
	AccountStatus string `json:"account_status"`
	Role          string `json:"-"`
}
