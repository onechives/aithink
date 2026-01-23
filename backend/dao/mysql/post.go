package mysql

import (
	"aithink/models"
	"errors"
	"fmt"
)

const (
	postColumns       = "p.id, p.title, p.summary, p.content_md, p.cover_url, p.category, p.tags, p.author_id, COALESCE(u.nickname, u.username) as author_name, p.status, p.reject_reason, p.like_count, p.created_at, p.updated_at"
	postSummaryFields = "p.id, p.title, p.summary, p.cover_url, p.category, p.tags, p.author_id, COALESCE(u.nickname, u.username) as author_name, p.status, p.reject_reason, p.like_count, p.created_at, p.updated_at"
)

func CreatePost(post *models.Post) error {
	// 新建文章记录
	query := "INSERT INTO posts (id, title, summary, content_md, cover_url, category, tags, author_id, status, like_count, created_at, updated_at) VALUES (:id, :title, :summary, :content_md, :cover_url, :category, :tags, :author_id, :status, :like_count, NOW(), NOW())"
	_, err := db.NamedExec(query, map[string]interface{}{
		"id":         post.ID,
		"title":      post.Title,
		"summary":    post.Summary,
		"content_md": post.Content,
		"cover_url":  post.CoverURL,
		"category":   post.Category,
		"tags":       post.Tags,
		"author_id":  post.AuthorID,
		"status":     post.Status,
		"like_count": post.LikeCount,
	})
	return err
}

// UpdatePost 更新文章字段（含审核状态与驳回原因）。
func UpdatePost(post *models.Post) error {
	query := "UPDATE posts SET title=?, summary=?, content_md=?, cover_url=?, category=?, tags=?, status=?, reject_reason=?, updated_at=NOW() WHERE id=?"
	result, err := db.Exec(query, post.Title, post.Summary, post.Content, post.CoverURL, post.Category, post.Tags, post.Status, post.RejectReason, post.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

// DeletePost 删除文章。
func DeletePost(id int64) error {
	result, err := db.Exec("DELETE FROM posts WHERE id=?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

// GetPostByID 按 ID 获取文章详情（含作者信息）。
func GetPostByID(id int64) (*models.Post, error) {
	post := new(models.Post)
	query := fmt.Sprintf("SELECT %s FROM posts p JOIN users u ON p.author_id = u.id WHERE p.id=?", postColumns)
	if err := db.Get(post, query, id); err != nil {
		return nil, err
	}
	return post, nil
}

// ListPosts 获取审核通过的文章列表，支持排序与搜索。
func ListPosts(sort, keyword string, page, size int) ([]models.PostSummary, int64, error) {
	items := make([]models.PostSummary, 0)
	offset := (page - 1) * size
	orderBy := "created_at DESC"
	if sort == "likes" {
		orderBy = "like_count DESC"
	}
	keywordLike := fmt.Sprintf("%%%s%%", keyword)
	var total int64
	if err := db.Get(&total, "SELECT COUNT(*) FROM posts WHERE status='approved' AND title LIKE ?", keywordLike); err != nil {
		return nil, 0, err
	}
	query := fmt.Sprintf("SELECT %s FROM posts p JOIN users u ON p.author_id = u.id WHERE p.status='approved' AND p.title LIKE ? ORDER BY %s LIMIT ? OFFSET ?", postSummaryFields, orderBy)
	if err := db.Select(&items, query, keywordLike, size, offset); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ListPostTitles 获取文章标题列表（侧边栏/搜索用）。
func ListPostTitles(sort, keyword string, size int) ([]models.PostTitle, error) {
	items := make([]models.PostTitle, 0)
	orderBy := "created_at DESC"
	if sort == "likes" {
		orderBy = "like_count DESC"
	}
	keywordLike := fmt.Sprintf("%%%s%%", keyword)
	query := fmt.Sprintf("SELECT id, title FROM posts WHERE status='approved' AND title LIKE ? ORDER BY %s LIMIT ?", orderBy)
	if err := db.Select(&items, query, keywordLike, size); err != nil {
		return nil, err
	}
	return items, nil
}

// IncrementLike 点赞并返回最新点赞数（仅审核通过的文章）。
func IncrementLike(id int64) (int64, error) {
	result, err := db.Exec("UPDATE posts SET like_count = like_count + 1 WHERE id=? AND status='approved'", id)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rows == 0 {
		return 0, errors.New("no rows affected")
	}
	var likeCount int64
	if err := db.Get(&likeCount, "SELECT like_count FROM posts WHERE id=?", id); err != nil {
		return 0, err
	}
	return likeCount, nil
}

// ListPostsByAuthor 获取作者自己的文章列表。
func ListPostsByAuthor(authorID int64, page, size int) ([]models.PostSummary, int64, error) {
	items := make([]models.PostSummary, 0)
	offset := (page - 1) * size
	var total int64
	if err := db.Get(&total, "SELECT COUNT(*) FROM posts WHERE author_id=?", authorID); err != nil {
		return nil, 0, err
	}
	query := fmt.Sprintf("SELECT %s FROM posts p JOIN users u ON p.author_id = u.id WHERE p.author_id=? ORDER BY p.created_at DESC LIMIT ? OFFSET ?", postSummaryFields)
	if err := db.Select(&items, query, authorID, size, offset); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ListPostsByStatus 按状态获取文章列表（管理端）。
func ListPostsByStatus(status string, page, size int) ([]models.PostSummary, int64, error) {
	items := make([]models.PostSummary, 0)
	offset := (page - 1) * size
	var total int64
	if err := db.Get(&total, "SELECT COUNT(*) FROM posts WHERE status=?", status); err != nil {
		return nil, 0, err
	}
	query := fmt.Sprintf("SELECT %s FROM posts p JOIN users u ON p.author_id = u.id WHERE p.status=? ORDER BY p.created_at DESC LIMIT ? OFFSET ?", postSummaryFields)
	if err := db.Select(&items, query, status, size, offset); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// UpdatePostStatus 更新文章审核状态（无驳回原因）。
func UpdatePostStatus(id int64, status string) error {
	result, err := db.Exec("UPDATE posts SET status=?, reject_reason='', updated_at=NOW() WHERE id=?", status, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

// UpdatePostStatusWithReason 更新文章审核状态并记录原因。
func UpdatePostStatusWithReason(id int64, status, reason string) error {
	result, err := db.Exec("UPDATE posts SET status=?, reject_reason=?, updated_at=NOW() WHERE id=?", status, reason, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
