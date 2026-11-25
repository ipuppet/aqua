// Package response writes JSON for HTTP-edge handlers.
// It maps domain errors (apperr). gRPC status conversion belongs in gateway/grpcx.
package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"aqua/pkg/apperr"
)

type errorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeSuccess(c *gin.Context, data any) {
	if data == nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, data)
}

// WriteError writes a standardized error payload with the given status and message.
func WriteError(c *gin.Context, status int, message string) {
	c.JSON(status, errorPayload{
		Code:    status,
		Message: message,
	})
}

// OK writes a successful response without payload.
func OK(c *gin.Context) {
	writeSuccess(c, nil)
}

// OKWithData writes a successful response with payload.
func OKWithData(c *gin.Context, data any) {
	writeSuccess(c, data)
}

// Fail maps the incoming error to a proper HTTP response and aborts the request.
func Fail(c *gin.Context, err error) {
	if err == nil {
		return
	}

	statusCode := http.StatusInternalServerError
	message := "internal server error"

	var mappedErr apperr.HTTPMappableError
	if errors.As(err, &mappedErr) {
		statusCode = mappedErr.HTTPStatus()
		message = mappedErr.ClientMessage()
	}

	c.AbortWithStatusJSON(statusCode, errorPayload{
		Code:    statusCode,
		Message: message,
	})
}

// OKOrFail writes OK when err is nil, otherwise maps err to a proper failure response.
func OKOrFail(c *gin.Context, err error) {
	if err != nil {
		Fail(c, err)
		return
	}

	OK(c)
}

// OKWithDataOrFail writes OKWithData when err is nil, otherwise maps err to a proper failure response.
func OKWithDataOrFail(c *gin.Context, data any, err error) {
	if err != nil {
		Fail(c, err)
		return
	}

	OKWithData(c, data)
}
