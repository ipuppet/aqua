package auth

import "time"

type Session struct {
	SessionID  string     `json:"session_id,omitempty"`
	UID        int        `json:"uid"`
	DeviceInfo string     `json:"device_info"`
	CreatedAt  *time.Time `json:"created_at"`
}
