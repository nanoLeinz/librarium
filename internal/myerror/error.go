package myerror

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/nanoLeinz/librarium/internal/model/dto"
)

type MyError struct {
	Code   int
	Status string
}

var (
	InternalServerErr = MyError{
		Code:   http.StatusInternalServerError,
		Status: "internal server error",
	}
)

func (s MyError) Error() string {
	return s.Status
}

func ToWebResponse(err MyError) *dto.WebResponse {

	log.Info(err)

	return &dto.WebResponse{
		Code:   err.Code,
		Status: err.Status,
		Result: nil,
	}
}

func NewNotFoundError(entity string) MyError {
	return MyError{
		Code:   http.StatusNotFound,
		Status: fmt.Sprintf("%s not found", entity),
	}
}

func NewDuplicateError(entity string) MyError {
	return MyError{
		Code:   http.StatusConflict,
		Status: fmt.Sprintf("%s already exist", entity),
	}
}

func NewBadRequestError(Status string) MyError {
	return MyError{
		Code:   http.StatusBadRequest,
		Status: Status,
	}
}
