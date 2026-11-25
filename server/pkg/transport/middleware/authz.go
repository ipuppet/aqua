package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ipuppet/aqua/server/pkg/apperr"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	"github.com/ipuppet/aqua/server/pkg/transport/authctx"
	"github.com/ipuppet/aqua/server/pkg/transport/response"
)

const allPermissionCode = "ALL"

type contextKeyPermissionSet struct{}

// rolePermissionReader is the RBAC read surface used by HTTP authorization.
// pkg/rbac.Reader satisfies this interface.
type rolePermissionReader interface {
	LoadRolePermissions(ctx context.Context, role string) ([]rbac.Permission, error)
}

// RequireAnyPermission checks that the subject has at least one required permission code.
func RequireAnyPermission(reader rolePermissionReader, codes ...string) gin.HandlerFunc {
	required := normalizeCodes(codes)

	return func(c *gin.Context) {
		if len(required) == 0 {
			c.Next()
			return
		}

		tokenUser := authctx.TokenUserFromContext(c.Request.Context())
		if tokenUser == nil {
			response.Fail(c, apperr.Unauthorized(""))
			return
		}

		permissionSet, err := getOrLoadPermissionSet(c, reader, tokenUser.Roles)
		if err != nil {
			response.Fail(c, apperr.Internal(""))
			return
		}
		if permissionSet[allPermissionCode] {
			c.Next()
			return
		}
		for _, code := range required {
			if permissionSet[code] {
				c.Next()
				return
			}
		}

		response.Fail(c, apperr.Forbidden("forbidden"))
	}
}

// RequireAllPermissions checks that the subject has all required permission codes.
func RequireAllPermissions(reader rolePermissionReader, codes ...string) gin.HandlerFunc {
	required := normalizeCodes(codes)

	return func(c *gin.Context) {
		if len(required) == 0 {
			c.Next()
			return
		}

		tokenUser := authctx.TokenUserFromContext(c.Request.Context())
		if tokenUser == nil {
			response.Fail(c, apperr.Unauthorized(""))
			return
		}

		permissionSet, err := getOrLoadPermissionSet(c, reader, tokenUser.Roles)
		if err != nil {
			response.Fail(c, apperr.Internal(""))
			return
		}
		if permissionSet[allPermissionCode] {
			c.Next()
			return
		}

		for _, code := range required {
			if !permissionSet[code] {
				response.Fail(c, apperr.Forbidden("forbidden"))
				return
			}
		}
		c.Next()
	}
}

// RequireAnyRole checks that the subject has at least one required role.
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	required := normalizeCodes(roles)

	return func(c *gin.Context) {
		if len(required) == 0 {
			c.Next()
			return
		}

		tokenUser := authctx.TokenUserFromContext(c.Request.Context())
		if tokenUser == nil {
			response.Fail(c, apperr.Unauthorized(""))
			return
		}

		roleSet := make(map[string]bool, len(tokenUser.Roles))
		for _, role := range tokenUser.Roles {
			if role == "" {
				continue
			}
			roleSet[role] = true
		}

		for _, role := range required {
			if roleSet[role] {
				c.Next()
				return
			}
		}

		response.Fail(c, apperr.Forbidden("forbidden"))
	}
}

func getOrLoadPermissionSet(c *gin.Context, reader rolePermissionReader, roles []string) (map[string]bool, error) {
	if cached, ok := c.Request.Context().Value(contextKeyPermissionSet{}).(map[string]bool); ok {
		return cached, nil
	}

	set := make(map[string]bool)
	for _, role := range roles {
		if role == "" {
			continue
		}
		perms, err := reader.LoadRolePermissions(c.Request.Context(), role)
		if err != nil {
			return nil, err
		}
		for _, p := range perms {
			if p.Code == "" {
				continue
			}
			set[p.Code] = true
		}
	}

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), contextKeyPermissionSet{}, set))
	return set, nil
}

func normalizeCodes(codes []string) []string {
	result := make([]string, 0, len(codes))
	seen := make(map[string]bool, len(codes))
	for _, code := range codes {
		code = strings.TrimSpace(code)
		if code == "" || seen[code] {
			continue
		}
		seen[code] = true
		result = append(result, code)
	}
	return result
}
