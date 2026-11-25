package httpbind

import (
	"github.com/gin-gonic/gin"

	"github.com/ipuppet/aqua/server/pkg/apperr"
	"github.com/ipuppet/aqua/server/pkg/transport/response"
)

// BindURI extracts URI parameters into obj.
func BindURI(c *gin.Context, obj any) bool {
	if err := c.BindUri(obj); err != nil {
		response.Fail(c, apperr.BadRequest("invalid URI parameters: "+err.Error()))
		return false
	}
	return true
}

// BindJSON unmarshals the JSON request body into obj.
func BindJSON(c *gin.Context, obj any) bool {
	if err := c.BindJSON(obj); err != nil {
		response.Fail(c, apperr.BadRequest("invalid JSON body: "+err.Error()))
		return false
	}
	return true
}

// Bind infers and binds request data into obj based on Content-Type.
func Bind(c *gin.Context, obj any) bool {
	if err := c.Bind(obj); err != nil {
		response.Fail(c, apperr.BadRequest("invalid request: "+err.Error()))
		return false
	}
	return true
}
