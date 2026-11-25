package auth

import "time"

type LoginAuditRecord struct {
	UID         int       `json:"uid"`
	LoginIP     string    `json:"login_ip"`
	DeviceInfo  string    `json:"device_info"`
	LoginMethod string    `json:"login_method"`
	RefreshJTI  string    `json:"refresh_jti"`
	CreatedAt   time.Time `json:"created_at"`
}
