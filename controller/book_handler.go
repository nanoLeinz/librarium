package controller

import (
	"encoding/json"
	"net/http"

	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/service"
	log "github.com/sirupsen/logrus"
)

type BookController struct {
	service service.BookService
	log     *log.Logger
}

func NewBookController(service service.BookService, log *log.Logger) *BookController {
	return &BookController{
		service: service,
		log:     log,
	}
}

func (s *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {

	rawUserData := r.Context().Value("memberDatas")

	userData := rawUserData.(map[string]any)

	s.log.WithFields(log.Fields{
		"function": "BookController.CreateBook",
		"memberId": userData["memberID"],
	}).Info("received request")

	Book := &dto.BookRequest{}

	if err := json.NewDecoder(r.Body).Decode(Book); err != nil {
		s.log.WithError(err).Error("Bad Request")

		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Result: nil,
		}

		helper.ResponseJSON(w, response)
		return
	}

	result, err := s.service.Create(r.Context(), Book)
	if err != nil {
		s.log.WithError(err).Error("failed to execute insert Book")

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
