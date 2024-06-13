package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppErrors types

type ErrorType int

const (
	ErrBadRequest      ErrorType = http.StatusBadRequest
	ErrUnauthorized    ErrorType = http.StatusUnauthorized
	ErrInternalFailure ErrorType = http.StatusInternalServerError
	ErrNotFound        ErrorType = http.StatusNotFound
)

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

/*
ErrorResponse and ErrorDetails help to create error json response of structure:
{
    "error": {
        "code": <err code>,
        "message": <err type>,
        "details": <err details>
    }
}
*/

// map const error values

var errorMessages = map[ErrorType]ErrorDetails{
	ErrBadRequest:      {Code: http.StatusBadRequest, Message: "Bad request"},
	ErrUnauthorized:    {Code: http.StatusUnauthorized, Message: "Authentication error"},
	ErrInternalFailure: {Code: http.StatusInternalServerError, Message: "Internal failure"},
	ErrNotFound:        {Code: http.StatusNotFound, Message: "Resource not found"},
}

func SendError(errType ErrorType, errDetails string, c *gin.Context) {

	errInfo, ok := errorMessages[errType]
	if !ok {
		// Handle unknown error
		errInfo = ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Unknown error",
		}
	}

	errInfo.Details = errDetails
	c.JSON(errInfo.Code, ErrorResponse{Error: errInfo})

}
