package models

import "time"

type NicknameRequest struct {
	ID           int64     `db:"id" json:"id"`
	UserID       int64     `db:"user_id" json:"userId"`
	Nickname     string    `db:"nickname" json:"nickname"`
	Status       string    `db:"status" json:"status"`
	RejectReason string    `db:"reject_reason" json:"rejectReason"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}
