package controller

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/service"
)

type AuthorController struct {
	service service.AuthorService
	log     *log.Logger
}

func NewAuthorController(service service.AuthorService, log *log.Logger) *AuthorController {
	return &AuthorController{
		service: service,
		log:     log,
	}
}

func (s *AuthorController) CreateAuthor(w http.ResponseWriter, r *http.Request) {

	rawUserData := r.Context().Value("memberDatas")

	userData := rawUserData.(map[string]any)

	s.log.WithFields(log.Fields{
		"function": "AuthorController.CreateAuthor",
		"memberId": userData["memberID"],
	}).Info("received request")

	author := &dto.AuthorRequest{}

	if err := json.NewDecoder(r.Body).Decode(author); err != nil {
		s.log.WithError(err).Error("Bad Request")

		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Result: nil,
		}

		helper.ResponseJSON(w, response)
		return
	}

	result, err := s.service.Create(r.Context(), author)
	if err != nil {
		s.log.WithError(err).Error("failed to execute insert author")

		response := myerror.ToWebResponse(err.(myerror.MyError))

		helper.ResponseJSON(w, response)
		return
	}

	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: result,
	}

	helper.ResponseJSON(w, response)
}
