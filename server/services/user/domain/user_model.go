package user

import (
	"time"

	"github.com/ipuppet/aqua/server/pkg/types"
)

// UserModelResponse represents the API response structure for a user record.
type UserModelResponse struct {
	UID         int        `json:"uid"`
	Username    string     `json:"username"`
	Nickname    string     `json:"nickname"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Avatar      string     `json:"avatar"`
	Gender      int        `json:"gender"`
	Lock        int        `json:"lock"`
	RegisterAt  *time.Time `json:"register_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

// UserModel represents a database record documenting a user's master information.
type UserModel struct {
	UID         int              `db:"uid" json:"uid"`
	Username    string           `db:"username" json:"username"`
	Nickname    string           `db:"nickname" json:"nickname"`
	Email       string           `db:"email" json:"email"`
	Phone       types.NullString `db:"phone" json:"phone"`
	Avatar      string           `db:"avatar" json:"avatar"`
	Gender      int              `db:"gender" json:"gender"`
	Password    string           `db:"password" json:"password"`
	Lock        int              `db:"lock" json:"lock"`
	RegisterAt  *time.Time       `db:"register_at" json:"register_at"`
	LastLoginAt *time.Time       `db:"last_login_at" json:"last_login_at"`
}
