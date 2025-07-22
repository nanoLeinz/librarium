package controller

import (

	// Add logging

	"github.com/go-playground/validator/v10"
	"github.com/nanoLeinz/librarium/service"
)

type MemberController struct {
	MemberService service.MemberService
	validator     *validator.Validate
}

func NewMemberController(service service.MemberService, validator *validator.Validate) *MemberController {
	return &MemberController{
		MemberService: service,
		validator:     validator,
	}
}
