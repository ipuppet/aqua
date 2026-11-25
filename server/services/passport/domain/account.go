package auth

import (
	"time"

	"aqua/pkg/types"
)

// AccountColumns is the rbac_user column list used by passport auth persistence.
const AccountColumns = "uid, username, nickname, email, phone, avatar, gender, password, lock, register_at, last_login_at"

// Account is the passport-owned credential row used for login, registration, and token claims.
type Account struct {
	UID         int              `db:"uid"`
	Username    string           `db:"username"`
	Nickname    string           `db:"nickname"`
	Email       string           `db:"email"`
	Phone       types.NullString `db:"phone"`
	Avatar      string           `db:"avatar"`
	Gender      int              `db:"gender"`
	Password    string           `db:"password"`
	Lock        int              `db:"lock"`
	RegisterAt  *time.Time       `db:"register_at"`
	LastLoginAt *time.Time       `db:"last_login_at"`
}
