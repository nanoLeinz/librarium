package myerror

import (
	"fmt"
	"net/http"

	"github.com/nanoLeinz/librarium/model/dto"
)

type MyError struct {
	Code   int
	Status string
}

var (
	InternalServerErr = &MyError{
		Code:   http.StatusInternalServerError,
		Status: "internal server error",
	}
)

func (s MyError) Error() string {
	return s.Status
}

func ToWebResponse(err MyError) *dto.WebResponse {
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
