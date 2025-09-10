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

	rawRq := dto.BookCopyRequest{}

	err := json.NewDecoder(r.Body).Decode(&rawRq)
	if err != nil {
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
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid book id",
			Result: nil,
		}
		helper.ResponseJSON(w, &response)
		return
	}

	err = s.copyService.Create(r.Context(), bookID, rawRq.Status, int(rawRq.Copies))
	if err != nil {
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, &response)
}

func (s *BookCopyController) UpdateStatus(w http.ResponseWriter, r *http.Request) {

	rawID := r.URL.Query().Get("copyID")

	copyID, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
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
		response := dto.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "invalid copy status",
			Result: nil,
		}

		helper.ResponseJSON(w, &response)
		return
	}

	err = s.copyService.Update(r.Context(), uint(copyID), &req)
	if err != nil {
		response := myerror.ToWebResponse(err.(myerror.MyError))
		helper.ResponseJSON(w, response)
		return
	}

	response := dto.WebResponse{
		Code:   http.StatusOK,
		Status: "success",
		Result: nil,
	}
	helper.ResponseJSON(w, &response)

}
