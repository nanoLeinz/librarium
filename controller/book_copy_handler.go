package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/service"
	log "github.com/sirupsen/logrus"
)

type BookCopyController struct {
	log         *log.Logger
	copyService service.BookCopyService
}

func NewBookCopyController(log *log.Logger, service service.BookCopyService) *BookCopyController {
	return &BookCopyController{
		log:         log,
		copyService: service,
	}
}

func (s *BookCopyController) logWithCtx(ctx context.Context, fun string) *log.Entry {
	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

	return s.log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": fun,
	})
}

func (s *BookCopyController) CreateCopies(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "BookCopyController.CreateCopies")

	rawRq := dto.BookCopyRequest{}
	err := json.NewDecoder(r.Body).Decode(&rawRq)
	if err != nil {
		logger.WithError(err).WithField("statusCode", http.StatusBadRequest).Error("invalid json request")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid json request",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	rawBookID := r.URL.Query().Get("bookID")
	bookID, err := uuid.Parse(rawBookID)
	if err != nil {
		logger.WithField("rawBookID", rawBookID).WithError(err).WithField("statusCode", http.StatusBadRequest).Error("invalid book id")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid book id",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	logger.WithFields(log.Fields{
		"bookID":     bookID,
		"status":     rawRq.Status,
		"copies":     rawRq.Copies,
		"statusCode": http.StatusOK,
	}).Info("received create book copies request")

	err = s.copyService.Create(r.Context(), bookID, rawRq.Status, int(rawRq.Copies))
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to create book copies")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"bookID":     bookID,
		"status":     rawRq.Status,
		"copies":     rawRq.Copies,
		"statusCode": http.StatusOK,
	}).Info("book copies created successfully")
	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, &response)
}

func (s *BookCopyController) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "BookCopyController.UpdateStatus")

	rawID := r.URL.Query().Get("copyID")
	copyID, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).WithField("statusCode", http.StatusBadRequest).Error("invalid copies")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid copies",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	req := dto.BookCopyRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Status == "" {
		logger.WithError(err).WithField("statusCode", http.StatusBadRequest).Error("invalid copy status")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid copy status",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	logger.WithFields(log.Fields{
		"copyID":     copyID,
		"status":     req.Status,
		"statusCode": http.StatusOK,
	}).Info("received update book copy status request")

	err = s.copyService.Update(r.Context(), uint(copyID), &req)
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to update book copy status")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"copyID":     copyID,
		"status":     req.Status,
		"statusCode": http.StatusOK,
	}).Info("book copy status updated successfully")
	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, &response)
}

func (s *BookCopyController) DeleteCopy(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "BookCopyController.DeleteCopy")

	rawID := r.URL.Query().Get("copyID")
	copyID, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).WithField("statusCode", http.StatusBadRequest).Error("invalid copies")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid copies",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	logger.WithFields(log.Fields{
		"copyID":     copyID,
		"statusCode": http.StatusOK,
	}).Info("received delete book copy request")

	err = s.copyService.DeleteById(r.Context(), uint(copyID))
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to delete book copy")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"copyID":     copyID,
		"statusCode": http.StatusOK,
	}).Info("book copy deleted successfully")
	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, &response)
}

func (s *BookCopyController) GetCopy(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "BookCopyController.GetCopy")

	rawID := r.URL.Query().Get("copyID")
	copyID, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).WithField("statusCode", http.StatusBadRequest).Error("invalid copies")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid copies",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	logger.WithFields(log.Fields{
		"copyID":     copyID,
		"statusCode": http.StatusOK,
	}).Info("received get book copy request")

	bookCopy, err := s.copyService.GetByID(r.Context(), uint(copyID))
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to get book copy")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"copyID":     copyID,
		"statusCode": http.StatusOK,
	}).Info("book copy fetched successfully")
	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: bookCopy,
	}
	helper.ResponseJSON(w, &response)
}

func (s *BookCopyController) GetAll(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "BookCopyController.GetAll")
	logger.WithField("statusCode", http.StatusOK).Info("received get all book copies request")

	copies, err := s.copyService.GetAll(r.Context())
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to get all book copies")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"count":      len(*copies),
		"statusCode": http.StatusOK,
	}).Info("all book copies fetched successfully")
	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: copies,
	}
	helper.ResponseJSON(w, &response)
}

func (s *BookCopyController) GetCopyByCondition(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "BookCopyController.GetCopyByCondition")

	req := dto.BookCopyRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.WithError(err).WithField("statusCode", http.StatusBadRequest).Error("invalid request")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid request",
			Result: nil,
		}

		helper.ResponseJSON(w, &response)
		return
	}

	logger.WithFields(log.Fields{
		"status":     req.Status,
		"bookID":     req.BookID,
		"statusCode": http.StatusOK,
	}).Info("received get book copies by condition request")

	copies, err := s.copyService.GetByCondition(r.Context(), &req)
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to get book copies by condition")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"count":      len(*copies),
		"statusCode": http.StatusOK,
	}).Info("book copies fetched by condition successfully")
	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: copies,
	}
	helper.ResponseJSON(w, &response)

}
