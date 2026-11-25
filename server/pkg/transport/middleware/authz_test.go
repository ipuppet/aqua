package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"aqua/pkg/rbac"
	"aqua/pkg/transport/authctx"
)

type rolePermissionStub struct {
	permissions []rbac.Permission
	err         error
	calls       int
}

func (s *rolePermissionStub) LoadRolePermissions(ctx context.Context, role string) ([]rbac.Permission, error) {
	s.calls++
	if s.err != nil {
		return nil, s.err
	}
	return s.permissions, nil
}

func TestRequireAnyPermission_ForbiddenOnMissingPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	reader := &rolePermissionStub{permissions: []rbac.Permission{{Code: "user.viewer"}}}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(authctx.WithTokenUser(c.Request.Context(), &authctx.TokenUser{UID: 7, Roles: []string{"user"}}))
		c.Next()
	})
	r.Use(RequireAnyPermission(reader, "user.manager"))
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestRequireAnyPermission_RequestScopeCache(t *testing.T) {
	gin.SetMode(gin.TestMode)
	reader := &rolePermissionStub{permissions: []rbac.Permission{{Code: "user.manager"}}}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(authctx.WithTokenUser(c.Request.Context(), &authctx.TokenUser{UID: 11, Roles: []string{"user"}}))
		c.Next()
	})
	r.Use(RequireAnyPermission(reader, "user.manager"))
	r.Use(RequireAnyPermission(reader, "user.manager"))
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
	if reader.calls != 1 {
		t.Fatalf("expected reader calls 1, got %d", reader.calls)
	}
}
