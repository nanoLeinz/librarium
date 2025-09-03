package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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

func (s *AuthorController) logWithCtx(ctx context.Context, fun string) *log.Entry {

	traceID := ctx.Value("traceID")

	return log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": fun,
	})
}

func (s *AuthorController) CreateAuthor(w http.ResponseWriter, r *http.Request) {

	logger := s.logWithCtx(r.Context(), "AuthorController.CreateAuthor")

	rawUserData := r.Context().Value("memberDatas")

	userData := rawUserData.(map[string]any)

	logger.WithField("memberId", userData["memberID"]).
		Info("received request")

	author := &dto.AuthorRequest{}

	if err := json.NewDecoder(r.Body).Decode(author); err != nil {
		logger.WithField("statusCode", http.StatusBadRequest).Error("Bad Request")

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
		logger.WithError(err).Error("failed to execute insert author")

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

func (s *AuthorController) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "AuthorController.UpdateAuthor")

	logger.Info()

	var rawID = r.PathValue("id")

	var authorID, err = strconv.ParseUint(rawID, 10, 8)

	if err != nil {
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Author ID not valid",
			Result: nil,
		}

		helper.ResponseJSON(w, response)
		return
	}

	var authorReq = dto.AuthorRequest{}

	decoder := json.NewDecoder(r.Body)

	decoder.Decode(&authorReq)

	err = s.service.Update(r.Context(), uint(authorID), &authorReq)

	if err != nil {
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

func (s *AuthorController) GetByID(w http.ResponseWriter, r *http.Request) {

	var rawID = r.PathValue("id")

	var authorID, err = strconv.ParseUint(rawID, 10, 8)

	if err != nil {
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Author ID not valid",
			Result: nil,
		}

		helper.ResponseJSON(w, response)
		return
	}

	rawRes, err := s.service.GetByIDs(r.Context(), uint(authorID))

	if err != nil {
		response := myerror.ToWebResponse(err.(myerror.MyError))

		helper.ResponseJSON(w, response)
		return
	}

	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: rawRes,
	}

	helper.ResponseJSON(w, response)
}

func (s *AuthorController) GetAllAuthor(w http.ResponseWriter, r *http.Request) {

	rawRes, err := s.service.GetAll(r.Context())

	if err != nil {
		response := myerror.ToWebResponse(err.(myerror.MyError))

		helper.ResponseJSON(w, response)
		return
	}

	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: rawRes,
	}

	helper.ResponseJSON(w, response)
}

func (s *AuthorController) GetAuthorsBook(w http.ResponseWriter, r *http.Request) {

	var rawID = r.PathValue("id")

	var authorID, err = strconv.ParseUint(rawID, 10, 8)

	if err != nil {
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Author ID not valid",
			Result: nil,
		}

		helper.ResponseJSON(w, response)
		return
	}

	var authorReq = dto.AuthorRequest{ID: uint(authorID)}

	rawRes, err := s.service.GetAuthorsBook(r.Context(), &authorReq)

	if err != nil {
		response := myerror.ToWebResponse(err.(myerror.MyError))

		helper.ResponseJSON(w, response)
		return
	}

	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: rawRes,
	}

	helper.ResponseJSON(w, response)
}

func (s *AuthorController) DeleteByID(w http.ResponseWriter, r *http.Request) {
	var rawID = r.PathValue("id")

	var authorID, err = strconv.ParseUint(rawID, 10, 8)

	if err != nil {
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Author ID not valid",
			Result: nil,
		}

		helper.ResponseJSON(w, response)
		return
	}

	err = s.service.DeleteById(r.Context(), uint(authorID))

	if err != nil {
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
