package models

import "time"

type Post struct {
	ID           int64     `db:"id" json:"id,string"`
	Title        string    `db:"title" json:"title"`
	Summary      string    `db:"summary" json:"summary"`
	Content      string    `db:"content_md" json:"content"`
	CoverURL     string    `db:"cover_url" json:"coverUrl"`
	Category     string    `db:"category" json:"category"`
	Tags         string    `db:"tags" json:"tags"`
	AuthorID     int64     `db:"author_id" json:"authorId"`
	Author       string    `db:"author_name" json:"author"`
	Status       string    `db:"status" json:"status"`
	RejectReason string    `db:"reject_reason" json:"rejectReason"`
	LikeCount    int64     `db:"like_count" json:"likeCount"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

type PostSummary struct {
	ID           int64     `db:"id" json:"id,string"`
	Title        string    `db:"title" json:"title"`
	Summary      string    `db:"summary" json:"summary"`
	CoverURL     string    `db:"cover_url" json:"coverUrl"`
	Category     string    `db:"category" json:"category"`
	Tags         string    `db:"tags" json:"tags"`
	AuthorID     int64     `db:"author_id" json:"authorId"`
	Author       string    `db:"author_name" json:"author"`
	Status       string    `db:"status" json:"status"`
	RejectReason string    `db:"reject_reason" json:"rejectReason"`
	LikeCount    int64     `db:"like_count" json:"likeCount"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

type PostTitle struct {
	ID    int64  `db:"id" json:"id,string"`
	Title string `db:"title" json:"title"`
}
