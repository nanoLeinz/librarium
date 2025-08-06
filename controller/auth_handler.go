package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/service"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	MemberService service.MemberService
	validator     *validator.Validate
	log           *logrus.Logger
}

func NewAuthController(service service.MemberService, validator *validator.Validate, log *logrus.Logger) *AuthController {
	return &AuthController{
		MemberService: service,
		validator:     validator,
		log:           log,
	}
}

func (s *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	s.log.WithField("function", "Register").Info("Received registration request")

	var req dto.MemberCreateRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		s.log.WithField("function", "Register").WithError(err).Warn("Failed to decode request body")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "bad request",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}
	s.log.WithField("function", "Register").Info("Request body decoded")

	req.AccountStatus = "ACTIVE"
	req.Role = "member"

	err := s.validator.Struct(&req)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"function": "Register",
			"email":    req.Email,
		}).WithError(err).Warn("Request validation failed")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "bad request",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}
	s.log.WithFields(logrus.Fields{
		"function": "Register",
		"email":    req.Email,
	}).Info("Request validation passed")

	member, err := s.MemberService.CreateMember(r.Context(), &req)

	if err != nil {
		s.log.WithFields(logrus.Fields{
			"function": "Register",
			"email":    req.Email,
		}).WithError(err).Error("Service error during member creation")

		response := myerror.ToWebResponse(err.(myerror.MyError))

		helper.ResponseJSON(w, response)
		return
	}

	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: *member,
	}
	helper.ResponseJSON(w, &response)

	s.log.WithFields(logrus.Fields{
		"function": "Register",
		"memberID": member.ID,
	}).Info("Registration successful and response sent")
}

func (s *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	s.log.WithField("function", "Login").Info("Received login request")

	req := &dto.LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.log.WithField("function", "Login").WithError(err).Warn("Failed to decode request body")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "bad request",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	if err := s.validator.Struct(req); err != nil {
		s.log.WithFields(logrus.Fields{
			"function": "Login",
			"email":    req.Email,
		}).WithError(err).Warn("Request validation failed")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "bad request",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}
	s.log.WithFields(logrus.Fields{
		"function": "Login",
		"email":    req.Email,
	}).Info("Request validation passed")

	member, err := s.MemberService.GetMemberByEmail(r.Context(), req.Email)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"function": "Login",
			"email":    req.Email,
		}).WithError(err).Warn("Authentication failed: member not found")

		response := myerror.ToWebResponse(err.(myerror.MyError))

		helper.ResponseJSON(w, response)
		return
	}

	ok := helper.CheckPassword(member.Password, req.Password)
	if !ok {
		s.log.WithFields(logrus.Fields{
			"function": "Login",
			"email":    req.Email,
			"memberID": member.ID,
		}).Warn("Authentication failed: wrong password")
		response := dto.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "unauthorized",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}
	s.log.WithFields(logrus.Fields{
		"function": "Login",
		"memberID": member.ID,
	}).Info("Password check successful")

	token, err := helper.GenerateJWTToken(member)

	if err != nil {
		s.log.WithFields(logrus.Fields{
			"function": "Login",
			"memberID": member.ID,
		}).WithError(err).Error("Failed to generate JWT token")
		response := dto.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "internal server error",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}
	s.log.WithFields(logrus.Fields{
		"function": "Login",
		"memberID": member.ID,
	}).Info("JWT token generated successfully")

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
	s.log.WithFields(logrus.Fields{
		"function": "Login",
		"memberID": member.ID,
	}).Info("Login successful and response sent")
}
