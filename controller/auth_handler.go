package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/service"
)

type AuthController struct {
	MemberService service.MemberService
	validator     *validator.Validate
}

func NewAuthController(service service.MemberService, validator *validator.Validate) *AuthController {
	return &AuthController{
		MemberService: service,
		validator:     validator,
	}
}

func (s *AuthController) Register(w http.ResponseWriter, r *http.Request) {

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

		helper.ResponseJSON(w, &response)

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

		helper.ResponseJSON(w, &response)

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

		helper.ResponseJSON(w, &response)

		return
	}

	log.Printf("CreateMember: member created: %+v", member) // Logging

	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: *member,
	}

	helper.ResponseJSON(w, &response)

	log.Println("CreateMember: response sent successfully")
}

func (s *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	// 1.Extract credentials (email, password) from request body.
	log.Println("Login Request Received")

	req := &dto.LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("Cannot Parse Request, error : %+v\n", err.Error())

		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "bad request",
			Result: nil,
		}

		helper.ResponseJSON(w, &response)

		log.Printf("Send Response with code : %d, rs : %+v\n", http.StatusBadRequest, &response)
		return
	}

	if err := s.validator.Struct(req); err != nil {
		log.Printf("Cannot Parse Request, error : %+v\n", err.Error())

		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "bad request",
			Result: nil,
		}

		helper.ResponseJSON(w, &response)

		log.Printf("Send Response with code : %d, rs : %+v\n", http.StatusBadRequest, &response)
		return
	}

	// 2. Find the user in the database by email.

	member, err := s.MemberService.GetMemberByEmail(r.Context(), req.Email)

	if err != nil {
		log.Println("Member not found")

		response := dto.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "unauthorized",
			Result: nil,
		}

		helper.ResponseJSON(w, &response)

		log.Printf("Send Response with code : %d, rs : %+v\n", http.StatusUnauthorized, &response)
		return
	}

	// 4. Compare the provided password with the stored hashed password.

	ok := helper.CheckPassword(member.Password, req.Password)

	if !ok {
		log.Println("wrong password")

		response := dto.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "unauthorized",
			Result: nil,
		}

		helper.ResponseJSON(w, &response)

		log.Printf("Send Response with code : %d, rs : %+v\n", http.StatusUnauthorized, &response)
		return
	}

	token, err := helper.GenerateJWTToken(member)

	if err != nil {
		log.Println("failed creating token")

		response := dto.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "failed creating token : " + err.Error(),
			Result: nil,
		}

		helper.ResponseJSON(w, &response)

		log.Printf("Send Response with code : %d, rs : %+v\n", http.StatusUnauthorized, &response)
		return
	}

	result := map[string]string{
		"id":    member.ID.String(),
		"email": member.Email,
		"token": token,
	}

	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "Success",
		Result: result,
	}

	helper.ResponseJSON(w, &response)
}
