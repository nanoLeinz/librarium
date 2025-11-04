package controller

import (

	// Add logging

	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"
	"github.com/nanoLeinz/librarium/internal/helper"
	"github.com/nanoLeinz/librarium/internal/model/dto"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/internal/service"
)

type MemberController struct {
	service   service.MemberService
	validator *validator.Validate
	log       *log.Logger
}

func NewMemberController(service service.MemberService, validator *validator.Validate, log *log.Logger) *MemberController {
	return &MemberController{
		service:   service,
		validator: validator,
		log:       log,
	}
}

func (s *MemberController) Profile(w http.ResponseWriter, r *http.Request) {

	memberDatas := r.Context().Value("memberDatas").(map[string]any)

	memberID := memberDatas["memberID"].(uuid.UUID)

	s.log.WithFields(log.Fields{
		"function": "member_handler.Profile",
		"memberID": memberID,
	}).Info("receive request Profile ")

	member, err := s.service.GetMemberByID(r.Context(), memberID)

	if err != nil {
		s.log.WithError(err).Error("failed to fetch member from database")

		response := myerror.ToWebResponse(err.(myerror.MyError))

		helper.ResponseJSON(w, response)
		return
	}

	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: member,
	}

	helper.ResponseJSON(w, response)

}

func (s *MemberController) DeleteProfile(w http.ResponseWriter, r *http.Request) {

	memberDatas := r.Context().Value("memberDatas").(map[string]any)

	memberID := memberDatas["memberID"].(uuid.UUID)

	s.log.WithFields(log.Fields{
		"function": "member_handler.DeleteProfile",
		"memberID": memberID,
	}).Info("receive request DeleteProfile ")

	err := s.service.DeleteMemberByID(r.Context(), memberID)

	if err != nil {
		s.log.WithError(err).Error("failed to delete member from database")

		response := myerror.ToWebResponse(err.(myerror.MyError))

		helper.ResponseJSON(w, response)
		return
	}

	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}

	helper.ResponseJSON(w, response)

}

func (s *MemberController) UpdateMember(w http.ResponseWriter, r *http.Request) {

	Data := r.Context().Value("memberDatas")

	MemberData := Data.(map[string]any)

	MemberID := MemberData["memberID"].(uuid.UUID)

	s.log.WithFields(log.Fields{
		"function": "Update Member",
		"MemberID": MemberID,
	}).Info("Received Update Member Request")

	var member *dto.MemberUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		s.log.WithError(err).Error("Bad Request")

		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Result: nil,
		}

		helper.ResponseJSON(w, response)
		return
	}

	member.ID = MemberID

	s.log.Info("gets data : ", member)

	err := s.service.UpdateMember(r.Context(), member)

	if err != nil {

		s.log.WithFields(log.Fields{
			"function": "Update Member",
			"memberID": MemberID,
		}).WithError(err).Error("Failed updating to database")

		response := myerror.ToWebResponse(err.(myerror.MyError))

		helper.ResponseJSON(w, response)
		return
	}

	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}

	helper.ResponseJSON(w, response)
}
