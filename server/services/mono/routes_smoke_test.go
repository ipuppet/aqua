package mono_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Ensures the path sets that mono merges onto one Gin engine do not panic on registration.
// Handlers are stubs; this only validates route tree compatibility across services.
func TestMergedHTTPRoutesDoNotConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	engine.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
		_, _ = c.Writer.Write([]byte("pang"))
	})

	mountAdmin := func(e *gin.Engine) {
		e.DELETE("/api/cache", stub)
		e.GET("/api/file/image/:module", stub)
		e.POST("/api/file/image/:module", stub)
		e.DELETE("/api/file/image/*path", stub)
	}
	mountPassport := func(e *gin.Engine) {
		e.Any("/api/login/password", stub)
		e.Any("/api/token", stub)
		e.Any("/api/register", stub)
		e.Any("/api/sessions", stub)
		e.Any("/api/sessions/*any", stub)
		e.Any("/api/users/:uid/sessions", stub)
		e.Any("/api/users/:uid/sessions/*any", stub)
	}
	mountUser := func(e *gin.Engine) {
		e.Any("/api/users/me", stub)
		e.Any("/api/users/me/*any", stub)
		e.Any("/api/users", stub)
		e.Any("/api/users/search/*any", stub)
		e.PUT("/api/users/:uid", stub)
		e.Any("/api/users/:uid/roles", stub)
		e.Any("/api/users/:uid/roles/*any", stub)
		e.Any("/api/users/:uid/permissions", stub)
		e.GET("/api/users/:uid", stub)
		e.Any("/api/access-control", stub)
		e.Any("/api/access-control/*any", stub)
	}
	mountBlog := func(e *gin.Engine) {
		e.GET("/api/articles", stub)
		e.GET("/api/articles/:aid", stub)
		e.POST("/api/articles", stub)
		e.GET("/api/articles/:aid/edit", stub)
		e.PUT("/api/articles/:aid", stub)
		e.DELETE("/api/articles/:aid", stub)
		e.GET("/api/comments/articles/:aid", stub)
		e.POST("/api/comments/articles/:aid", stub)
		e.GET("/api/comments/:id/children", stub)
		e.GET("/api/about", stub)
		e.PUT("/api/about", stub)
	}
	require.NotPanics(t, func() {
		mountAdmin(engine)
		mountPassport(engine)
		mountUser(engine)
		mountBlog(engine)
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "pang", w.Body.String())
}

func stub(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
