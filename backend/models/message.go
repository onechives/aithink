package models

import "time"

const (
	MessageUnread = "unread"
	MessageRead   = "read"
)

type Message struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"userId"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
