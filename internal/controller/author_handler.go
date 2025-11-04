package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/nanoLeinz/librarium/internal/helper"
	"github.com/nanoLeinz/librarium/internal/model/dto"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/internal/service"
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

	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

	return s.log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": fun,
	})
}

func (s *AuthorController) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "AuthorController.CreateAuthor")

	rawUserData := r.Context().Value("memberDatas")
	userData := rawUserData.(map[string]any)

	logger.WithField("memberId", userData["memberID"]).Info("received create author request")

	author := &dto.AuthorRequest{}
	if err := json.NewDecoder(r.Body).Decode(author); err != nil {
		logger.WithField("statusCode", http.StatusBadRequest).Error("bad request: failed to decode body")
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

	logger.WithField("authorName", author.Name).Info("author created successfully")
	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: result,
	}
	helper.ResponseJSON(w, response)
}

func (s *AuthorController) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "AuthorController.UpdateAuthor")
	rawID := r.PathValue("id")
	authorID, err := strconv.ParseUint(rawID, 10, 32)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).Error("invalid author ID")
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Author ID not valid",
			Result: nil,
		}
		helper.ResponseJSON(w, response)
		return
	}

	var authorReq dto.AuthorRequest
	if err := json.NewDecoder(r.Body).Decode(&authorReq); err != nil {
		logger.WithField("statusCode", http.StatusBadRequest).Error("bad request: failed to decode body")
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Result: nil,
		}
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithFields(log.Fields{
		"authorID":   authorID,
		"authorName": authorReq.Name,
	}).Info("received update author request")

	err = s.service.Update(r.Context(), uint(authorID), &authorReq)
	if err != nil {
		logger.WithError(err).Error("failed to update author")
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithField("authorID", authorID).Info("author updated successfully")
	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, response)
}

func (s *AuthorController) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "AuthorController.GetByID")
	rawID := r.PathValue("id")
	authorID, err := strconv.ParseUint(rawID, 10, 8)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).Error("invalid author ID")
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Author ID not valid",
			Result: nil,
		}
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithField("authorID", authorID).Info("received get author by ID request")
	rawRes, err := s.service.GetByIDs(r.Context(), uint(authorID))
	if err != nil {
		logger.WithError(err).Error("failed to get author by ID")
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithField("authorID", authorID).Info("author fetched successfully")
	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: rawRes,
	}
	helper.ResponseJSON(w, response)
}

func (s *AuthorController) GetAllAuthor(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "AuthorController.GetAllAuthor")
	logger.Info("received get all authors request")

	rawRes, err := s.service.GetAll(r.Context())
	if err != nil {
		logger.WithError(err).Error("failed to get all authors")
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithField("count", len(*rawRes)).Info("all authors fetched successfully")
	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: rawRes,
	}
	helper.ResponseJSON(w, response)
}

func (s *AuthorController) GetAuthorsBook(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "AuthorController.GetAuthorsBook")
	rawID := r.PathValue("id")
	authorID, err := strconv.ParseUint(rawID, 10, 8)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).Error("invalid author ID")
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Author ID not valid",
			Result: nil,
		}
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithField("authorID", authorID).Info("received get author's books request")
	authorReq := dto.AuthorRequest{ID: uint(authorID)}
	rawRes, err := s.service.GetAuthorsBook(r.Context(), &authorReq)
	if err != nil {
		logger.WithError(err).Error("failed to get author's books")
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithField("authorID", authorID).Info("author's books fetched successfully")
	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: rawRes,
	}
	helper.ResponseJSON(w, response)
}

func (s *AuthorController) DeleteByID(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "AuthorController.DeleteByID")
	rawID := r.PathValue("id")
	authorID, err := strconv.ParseUint(rawID, 10, 8)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).Error("invalid author ID")
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Author ID not valid",
			Result: nil,
		}
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithField("authorID", authorID).Info("received delete author request")
	err = s.service.DeleteById(r.Context(), uint(authorID))
	if err != nil {
		logger.WithError(err).Error("failed to delete author")
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithField("authorID", authorID).Info("author deleted successfully")
	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, response)
}
