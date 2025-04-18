package contract

import (
	"net/http"
	"pie_fire_dine/errs"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type Response struct {
	StatusCode int             `json:"-"`
	Errors     []ErrorResponse `json:"errors,omitempty"`
	Data       interface{}     `json:"data"`
}

func NewResponse(status int, data interface{}) Response {
	return Response{
		StatusCode: status,
		Data:       data,
	}
}

func NewErrorResponse(err error) Response {
	var status int
	var errors []ErrorResponse
	switch errorRes := err.(type) {
	case errs.AppError:
		status = errorRes.Code
		errResponse := ErrorResponse{
			Message: errorRes.Message,
		}
		errors = append(errors, errResponse)
	default:
		status = http.StatusInternalServerError
		errResponse := ErrorResponse{
			Message: "Internal Server Error",
		}
		errors = append(errors, errResponse)
	}
	return Response{
		StatusCode: status,
		Errors:     errors,
	}
}
