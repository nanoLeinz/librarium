package controller

import (
	"encoding/json"
	"log" // Add logging
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/model/dto"
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

func (s MemberController) CreateMember(w http.ResponseWriter, r *http.Request) {

	log.Println("CreateMember: received request") // Logging

	var req dto.MemberCreateRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateMember: failed to decode request body: %v", err) // Logging
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "bad request",
			Result: nil,
		}

		helper.ResponseJSON(w, http.StatusBadRequest, response)

		return
	}

	log.Printf("CreateMember: request decoded: %+v", req) // Logging

	req.AccountStatus = "ACTIVE"
	log.Println("CreateMember: set AccountStatus to ACTIVE") // Logging

	req.Role = "member"
	log.Println("CreateMember: set Role to Member")

	err := s.validator.Struct(&req)

	if err != nil {
		log.Printf("CreateMember: validation failed: %v", err) // Logging
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "bad request",
			Result: nil,
		}

		helper.ResponseJSON(w, http.StatusBadRequest, response)

		return
	}

	log.Println("CreateMember: validation passed") // Logging

	member, err := s.MemberService.CreateMember(r.Context(), &req)

	if err != nil {
		log.Printf("CreateMember: service error: %v", err) // Logging
		response := dto.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: err.Error(),
			Result: nil,
		}

		helper.ResponseJSON(w, http.StatusInternalServerError, response)

		return
	}

	log.Printf("CreateMember: member created: %+v", member) // Logging

	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: *member,
	}

	helper.ResponseJSON(w, http.StatusOK, response)

	log.Println("CreateMember: response sent successfully")
}
