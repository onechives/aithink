package logic

import (
	"aithink/dao/mysql"
	"aithink/models"
	"aithink/pkg/snowflake"
	"errors"
)

// ErrNotFound 未找到文章
var ErrNotFound = errors.New("post not found")

// CreatePost 创建文章（生成雪花 ID）。
func CreatePost(post *models.Post) error {
	post.ID = snowflake.GenIDByInt()
	return mysql.CreatePost(post)
}

// UpdatePost 更新文章内容。
func UpdatePost(post *models.Post) error {
	return mysql.UpdatePost(post)
}

// DeletePost 删除文章。
func DeletePost(id int64) error {
	return mysql.DeletePost(id)
}

// GetPostDetail 获取文章详情（若不存在返回 ErrNotFound）。
func GetPostDetail(id int64) (*models.Post, error) {
	post, err := mysql.GetPostByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return post, nil
}

// ListPosts 获取已审核通过的文章列表。
func ListPosts(sort, keyword string, page, size int) ([]models.PostSummary, int64, error) {
	return mysql.ListPosts(sort, keyword, page, size)
}

// ListPostsByAuthor 获取指定作者的文章列表。
func ListPostsByAuthor(authorID int64, page, size int) ([]models.PostSummary, int64, error) {
	return mysql.ListPostsByAuthor(authorID, page, size)
}

// ListPostsByStatus 按状态获取文章列表（管理端使用）。
func ListPostsByStatus(status string, page, size int) ([]models.PostSummary, int64, error) {
	return mysql.ListPostsByStatus(status, page, size)
}

// UpdatePostStatus 更新文章审核状态，并发送站内信通知。
func UpdatePostStatus(id int64, status string) error {
	post, err := mysql.GetPostByID(id)
	if err != nil {
		return err
	}
	if err := mysql.UpdatePostStatus(id, status); err != nil {
		return err
	}
	if status == models.StatusApproved {
		return SendMessage(post.AuthorID, "文章审核通过", "你的文章已审核通过并公开展示。")
	}
	if status == models.StatusRejected {
		return SendMessage(post.AuthorID, "文章审核未通过", "你的文章审核未通过，仅作者可见。")
	}
	return nil
}

// UpdatePostStatusWithReason 更新文章审核状态并附带原因。
func UpdatePostStatusWithReason(id int64, status, reason string) error {
	post, err := mysql.GetPostByID(id)
	if err != nil {
		return err
	}
	if err := mysql.UpdatePostStatusWithReason(id, status, reason); err != nil {
		return err
	}
	if status == models.StatusApproved {
		return SendMessage(post.AuthorID, "文章审核通过", "你的文章已审核通过并公开展示。")
	}
	if status == models.StatusRejected {
		content := "你的文章审核未通过，仅作者可见。"
		if reason != "" {
			content = content + " 原因：" + reason
		}
		return SendMessage(post.AuthorID, "文章审核未通过", content)
	}
	return nil
}

// ListPostTitles 获取文章标题列表（侧边栏/搜索）。
func ListPostTitles(sort, keyword string, size int) ([]models.PostTitle, error) {
	return mysql.ListPostTitles(sort, keyword, size)
}

// LikePost 点赞文章。
func LikePost(id int64) (int64, error) {
	return mysql.IncrementLike(id)
}
