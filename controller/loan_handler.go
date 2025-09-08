package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/internal/myerror"
	"github.com/nanoLeinz/librarium/model/dto"
	"github.com/nanoLeinz/librarium/service"
	log "github.com/sirupsen/logrus"
)

type LoanController struct {
	log     *log.Logger
	service service.LoanService
}

func NewLoanController(log *log.Logger, service service.LoanService) *LoanController {
	return &LoanController{
		log:     log,
		service: service,
	}
}

func (s *LoanController) logWithCtx(ctx context.Context, fun string) *log.Entry {
	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

	return s.log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": fun,
	})
}

func (s *LoanController) CreateLoan(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "LoanController.CreateLoan")

	rawReq := dto.LoanRequest{}
	err := json.NewDecoder(r.Body).Decode(&rawReq)
	if err != nil {
		logger.WithError(err).WithField("statusCode", http.StatusBadRequest).Error("invalid request: failed to decode body")
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid request",
			Result: nil,
		}
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithFields(log.Fields{
		"memberID":   rawReq.MemberID,
		"bookCopyID": rawReq.BookCopyID,
		"statusCode": http.StatusOK,
	}).Info("received create loan request")

	res, err := s.service.Create(r.Context(), &rawReq)
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to create loan")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"memberID":   rawReq.MemberID,
		"bookCopyID": rawReq.BookCopyID,
		"statusCode": http.StatusOK,
	}).Info("loan created successfully")
	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: res,
	}
	helper.ResponseJSON(w, &response)
}

func (s *LoanController) UpdateLoan(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "LoanController.UpdateLoan")

	rawID := r.URL.Query().Get("id")
	loanID, err := uuid.Parse(rawID)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).WithField("statusCode", http.StatusBadRequest).Error("invalid loan id")
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid loan id",
			Result: nil,
		}
		helper.ResponseJSON(w, response)
		return
	}

	rawReq := dto.LoanRequest{}
	err = json.NewDecoder(r.Body).Decode(&rawReq)
	if err != nil {
		logger.WithError(err).WithField("statusCode", http.StatusBadRequest).Error("invalid request: failed to decode body")
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid request",
			Result: nil,
		}
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithFields(log.Fields{
		"loanID":     loanID,
		"statusCode": http.StatusOK,
	}).Info("received update loan request")

	err = s.service.Update(r.Context(), loanID, &rawReq)
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to update loan")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"loanID":     loanID,
		"statusCode": http.StatusOK,
	}).Info("loan updated successfully")
	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, &response)
}

func (s *LoanController) DeleteLoan(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "LoanController.DeleteLoan")

	rawID := r.URL.Query().Get("id")
	loanID, err := uuid.Parse(rawID)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).WithField("statusCode", http.StatusBadRequest).Error("invalid loan id")
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid loan id",
			Result: nil,
		}
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithFields(log.Fields{
		"loanID":     loanID,
		"statusCode": http.StatusOK,
	}).Info("received delete loan request")

	err = s.service.DeleteById(r.Context(), loanID)
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to delete loan")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"loanID":     loanID,
		"statusCode": http.StatusOK,
	}).Info("loan deleted successfully")
	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, &response)
}

func (s *LoanController) GetLoanByID(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "LoanController.GetLoanByID")

	rawID := r.URL.Query().Get("id")
	loanID, err := uuid.Parse(rawID)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).WithField("statusCode", http.StatusBadRequest).Error("invalid loan id")
		response := &dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid loan id",
			Result: nil,
		}
		helper.ResponseJSON(w, response)
		return
	}

	logger.WithFields(log.Fields{
		"loanID":     loanID,
		"statusCode": http.StatusOK,
	}).Info("received get loan by ID request")

	res, err := s.service.GetByID(r.Context(), loanID)
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to get loan by ID")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"loanID":     loanID,
		"statusCode": http.StatusOK,
	}).Info("loan fetched successfully")
	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: res,
	}
	helper.ResponseJSON(w, &response)
}

func (s *LoanController) GetAllLoan(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "LoanController.GetAllLoan")
	logger.WithField("statusCode", http.StatusOK).Info("received get all loans request")

	loans, err := s.service.GetAll(r.Context())
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to get all loans")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"count":      len(*loans),
		"statusCode": http.StatusOK,
	}).Info("all loans fetched successfully")
	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: loans,
	}
	helper.ResponseJSON(w, &response)
}
