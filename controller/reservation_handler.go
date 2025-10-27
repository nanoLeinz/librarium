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

type ReservationController struct {
	log  *log.Logger
	serv service.ReservationService
}

func NewReservationController(log *log.Logger, service service.ReservationService) *ReservationController {
	return &ReservationController{
		log:  log,
		serv: service,
	}
}

func (s *ReservationController) logWithCtx(ctx context.Context, fun string) *log.Entry {

	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

	return s.log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": fun,
	})
}

func (s *ReservationController) CreateReservation(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "Controller.CreateReservation")

	rawRq := dto.ReservationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&rawRq); err != nil {
		logger.WithError(err).WithField("statusCode", http.StatusBadRequest).Error("failed to decode create reservation request")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	logger.WithFields(log.Fields{
		"bookID":     rawRq.BookID,
		"memberID":   rawRq.MemberID,
		"statusCode": http.StatusOK,
	}).Info("received create reservation request")

	res, err := s.serv.Create(r.Context(), &rawRq)
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to create reservation")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"reservationID": res.ID,
		"queuePosition": res.QueuePosition,
		"statusCode":    http.StatusOK,
	}).Info("reservation created successfully")

	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: res,
	}
	helper.ResponseJSON(w, &response)
}

func (s *ReservationController) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "Controller.UpdateReservation")

	rawRq := dto.ReservationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&rawRq); err != nil {
		logger.WithError(err).WithField("statusCode", http.StatusBadRequest).Error("failed to decode update reservation request")
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	rawID := r.PathValue("id")
	ID, err := uuid.Parse(rawID)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).WithField("statusCode", http.StatusNotFound).Error("invalid reservation id")
		response := dto.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Reservation Not Found",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	logger.WithFields(log.Fields{
		"reservationID": ID,
		"bookID":        rawRq.BookID,
		"memberID":      rawRq.MemberID,
		"statusCode":    http.StatusOK,
	}).Info("received update reservation request")

	if err := s.serv.Update(r.Context(), ID, &rawRq); err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to update reservation")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"reservationID": ID,
		"statusCode":    http.StatusOK,
	}).Info("reservation updated successfully")

	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, &response)
}

func (s *ReservationController) DeleteReservation(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "Controller.DeleteReservation")

	rawID := r.PathValue("id")
	ID, err := uuid.Parse(rawID)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).WithField("statusCode", http.StatusNotFound).Error("invalid reservation id")
		response := dto.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Reservation Not Found",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	logger.WithFields(log.Fields{
		"reservationID": ID,
		"statusCode":    http.StatusOK,
	}).Info("received delete reservation request")

	if err := s.serv.DeleteById(r.Context(), ID); err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to delete reservation")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"reservationID": ID,
		"statusCode":    http.StatusOK,
	}).Info("reservation deleted successfully")

	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, &response)
}

func (s *ReservationController) GetReservationByID(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "Controller.GetReservationByID")

	rawID := r.PathValue("id")
	ID, err := uuid.Parse(rawID)
	if err != nil {
		logger.WithField("rawID", rawID).WithError(err).WithField("statusCode", http.StatusNotFound).Error("invalid reservation id")
		response := dto.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Reservation Not Found",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	logger.WithFields(log.Fields{
		"reservationID": ID,
		"statusCode":    http.StatusOK,
	}).Info("received get reservation by ID request")

	res, err := s.serv.GetByID(r.Context(), ID)
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to get reservation by ID")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"reservationID": res.ID,
		"statusCode":    http.StatusOK,
	}).Info("reservation fetched successfully")

	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: res,
	}
	helper.ResponseJSON(w, &response)
}

func (s *ReservationController) GetAllReservation(w http.ResponseWriter, r *http.Request) {
	logger := s.logWithCtx(r.Context(), "Controller.GetAllReservation")
	logger.WithField("statusCode", http.StatusOK).Info("received get all reservations request")

	res, err := s.serv.GetAll(r.Context())
	if err != nil {
		webRes := myerror.ToWebResponse(err.(myerror.MyError))
		logger.WithError(err).WithField("statusCode", webRes.Code).Error("failed to get all reservations")
		helper.ResponseJSON(w, webRes)
		return
	}

	logger.WithFields(log.Fields{
		"count":      len(res),
		"statusCode": http.StatusOK,
	}).Info("all reservations fetched successfully")

	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: res,
	}
	helper.ResponseJSON(w, &response)
}
