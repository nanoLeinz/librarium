package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/internal/helper"
	"github.com/nanoLeinz/librarium/internal/model/dto"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/internal/service"
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

func (s *BookController) logWithCtx(ctx context.Context, fun string) *log.Entry {

	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

	return s.log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": fun,
	})
}

func (s *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "BookController.CreateBook")

	rawUserData := r.Context().Value("memberDatas")
	userData := rawUserData.(map[string]any)
	logger.WithField("memberId", userData["memberID"]).Info("received create book request")

	Book := &dto.BookRequest{}
	if err := json.NewDecoder(r.Body).Decode(Book); err != nil {
		logger.WithError(err).WithField("statusCode", http.StatusBadRequest).
			Error("bad request: failed to decode body")
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
		logger.WithError(err).
			WithField("statusCode", err.(myerror.MyError).Code).
			Error("failed to execute insert book")
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithFields(log.Fields{
		"title":      Book.Title,
		"statusCode": http.StatusOK,
	}).Info("book created successfully")
	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: result,
	}
	helper.ResponseJSON(w, response)
}

func (s *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "BookController.UpdateBook")

	rawID := r.PathValue("id")
	bookID, err := uuid.Parse(rawID)
	if err != nil {
		logger.WithField("rawID", rawID).
			WithError(err).
			WithField("statusCode", http.StatusBadRequest).
			Error("invalid book id")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid book id",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	bookReq := dto.BookRequest{}
	if err := json.NewDecoder(r.Body).Decode(&bookReq); err != nil {
		logger.WithError(err).
			WithField("statusCode", http.StatusBadRequest).
			Error("bad request: failed to decode body")
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Result: nil,
		}
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithFields(log.Fields{
		"bookID":     bookID,
		"title":      bookReq.Title,
		"statusCode": http.StatusOK,
	}).Info("received update book request")

	err = s.service.Update(r.Context(), bookID, &bookReq)
	if err != nil {
		logger.WithError(err).
			WithField("statusCode", err.(myerror.MyError).Code).
			Error("failed to update book")
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithFields(log.Fields{
		"bookID":     bookID,
		"statusCode": http.StatusOK,
	}).Info("book updated successfully")
	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, response)
}

func (s *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "BookController.DeleteBook")

	rawID := r.PathValue("id")
	bookID, err := uuid.Parse(rawID)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).Error("invalid book id")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid book id",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	logger.WithField("bookID", bookID).Info("received delete book request")
	err = s.service.DeleteByID(r.Context(), bookID)
	if err != nil {
		logger.WithError(err).Error("failed to delete book")
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithField("bookID", bookID).Info("book deleted successfully")
	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, response)
}

func (s *BookController) GetAll(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "BookController.GetAll")
	logger.Info("received get all books request")

	books, err := s.service.GetAll(r.Context())
	if err != nil {
		logger.WithError(err).Error("failed to get all books")
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithField("count", len(*books)).Info("all books fetched successfully")
	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: books,
	}
	helper.ResponseJSON(w, response)
}

func (s *BookController) GetBook(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "BookController.GetBook")

	rawID := r.PathValue("id")
	bookID, err := uuid.Parse(rawID)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).Error("invalid book id")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid book id",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	logger.WithField("bookID", bookID).Info("received get book by ID request")
	book, err := s.service.GetByID(r.Context(), bookID)
	if err != nil {
		logger.WithError(err).Error("failed to get book by ID")
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithField("bookID", bookID).Info("book fetched successfully")
	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: book,
	}
	helper.ResponseJSON(w, response)
}

func (s *BookController) GetBookByTitle(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "BookController.GetBookByTitle")

	title := r.URL.Query().Get("s")
	logger.WithField("title", title).Info("received get book by title request")
	books, err := s.service.GetByTitle(r.Context(), title)
	if err != nil {
		logger.WithError(err).Error("failed to get books by title")
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithField("count", len(*books)).Info("books fetched by title successfully")
	response := &dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: books,
	}
	helper.ResponseJSON(w, response)
}
