package rest_errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

//RestError
type RestError interface {
	Message() string
	Status() int64
	Code() string
}

type restError struct {
	RMessage string `json:"message"`
	RStatus  int64  `json:"status"`
	RCode    string `json:"code"`
}

func (re *restError) Status() int64 {
	return re.RStatus
}

func (re *restError) Message() string {
	return re.RMessage
}

func (re *restError) Code() string {
	return re.RCode
}

//NewRestErrorFromBytes : creates the resterror domain by taking the paramter as a slice of bytes
func NewRestErrorFromBytes(byte []byte) (RestError, error) {
	var apiErr restError
	if err := json.Unmarshal(byte, &apiErr); err != nil {
		return nil, errors.New("Invalid json")
	}
	return &apiErr, nil
}

//NewBadRequestError : this method implements the bad request error
func NewBadRequestError(message string) RestError {
	return &restError{
		RMessage: message,
		RStatus:  http.StatusBadRequest,
		RCode:    "bad_request",
	}
}

//NewNotFoundError : this method implements the not found error
func NewNotFoundError(message string) RestError {
	return &restError{
		RMessage: message,
		RStatus:  http.StatusNotFound,
		RCode:    "not_found",
	}
}

//NewInternalServerError : this method implements the internal server error
func NewInternalServerError(message string) RestError {
	return &restError{
		RMessage: message,
		RStatus:  http.StatusInternalServerError,
		RCode:    "internal_server_error",
	}
}

//NewUnauthorizedError : this method implements the unauthorized error
func NewUnauthorizedError(message string) RestError {
	return &restError{
		RMessage: message,
		RStatus:  http.StatusUnauthorized,
		RCode:    "unauthorized_error",
	}
}
