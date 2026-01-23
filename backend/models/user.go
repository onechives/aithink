package models

import "time"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"

	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"

	NicknamePending  = "pending"
	NicknameApproved = "approved"
	NicknameRejected = "rejected"
)

type User struct {
	ID           int64     `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Nickname     string    `db:"nickname" json:"nickname"`
	Password     string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	Status       string    `db:"status" json:"status"`
	RejectReason string    `db:"reject_reason" json:"rejectReason"`
	TOTPSecret   string    `db:"totp_secret" json:"-"`
	TOTPEnabled  bool      `db:"totp_enabled" json:"totpEnabled"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}
