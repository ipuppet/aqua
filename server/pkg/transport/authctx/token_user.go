package authctx

import "context"

type contextKeyTokenUser struct{}

// ContextKeyTokenUser is the key for storing *TokenUser in context.Context.
var ContextKeyTokenUser = contextKeyTokenUser{}

// WithTokenUser stores token user in context and returns a derived context.
func WithTokenUser(ctx context.Context, user *TokenUser) context.Context {
	return context.WithValue(ctx, ContextKeyTokenUser, user)
}

// TokenUserFromContext returns the TokenUser from context, or nil.
func TokenUserFromContext(ctx context.Context) *TokenUser {
	u, _ := ctx.Value(ContextKeyTokenUser).(*TokenUser)
	return u
}

// TokenUser is the cross-service identity projection extracted from JWT claims.
// Services must not import another service's UserModel or Account for request auth.
type TokenUser struct {
	UID   int
	Roles []string
}
