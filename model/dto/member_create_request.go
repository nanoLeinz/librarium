package dto

type MemberCreateRequest struct {
	Email         string `validate:"required,email"`
	Password      string `validate:"required"`
	FullName      string `validate:"required"`
	AccountStatus string `validate:"required"`
}
