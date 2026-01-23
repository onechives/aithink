package controllers

import (
	"aithink/logic"
	"aithink/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

type PostCreateRequest struct {
	// 文章标题
	Title    string `json:"title" binding:"required"`
	// 文章摘要
	Summary  string `json:"summary" binding:"required"`
	// Markdown 正文
	Content  string `json:"content" binding:"required"`
	// 封面图地址
	CoverURL string `json:"coverUrl"`
	// 分类
	Category string `json:"category" binding:"required"`
	// 标签（逗号分隔）
	Tags     string `json:"tags"`
}

type PostUpdateRequest struct {
	// 更新用字段，与创建保持一致
	Title    string `json:"title" binding:"required"`
	Summary  string `json:"summary" binding:"required"`
	Content  string `json:"content" binding:"required"`
	CoverURL string `json:"coverUrl"`
	Category string `json:"category" binding:"required"`
	Tags     string `json:"tags"`
}

type PostListResponse struct {
	// 文章列表 + 总数
	Items []models.PostSummary `json:"items"`
	Total int64                `json:"total"`
}

// PostCreateHandler 创建文章；普通用户为待审核状态。
func PostCreateHandler(c *gin.Context) {
	userID := c.GetInt64(ContextUserIDKey)
	if userID == 0 {
		ResponseError(c, CodeNeedLogin)
		return
	}
	role := c.GetString(ContextUserRoleKey)
	var req PostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	post := &models.Post{
		Title:    req.Title,
		Summary:  req.Summary,
		Content:  req.Content,
		CoverURL: req.CoverURL,
		Category: req.Category,
		Tags:     req.Tags,
		AuthorID: userID,
		Status:   models.StatusPending,
	}
	if role == models.RoleAdmin {
		post.Status = models.StatusApproved
	}
	if err := logic.CreatePost(post); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"id": strconv.FormatInt(post.ID, 10)})
}

// PostUpdateHandler 更新文章内容；非管理员更新后重新进入待审核。
func PostUpdateHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	userID := c.GetInt64(ContextUserIDKey)
	role := c.GetString(ContextUserRoleKey)
	var req PostUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	existing, err := logic.GetPostDetail(id)
	if err != nil {
		ResponseError(c, CodeNotFound)
		return
	}
	if role != models.RoleAdmin && existing.AuthorID != userID {
		ResponseError(c, CodeAuthFailed)
		return
	}
	status := existing.Status
	if role != models.RoleAdmin {
		status = models.StatusPending
	}
	post := &models.Post{
		ID:           id,
		Title:        req.Title,
		Summary:      req.Summary,
		Content:      req.Content,
		CoverURL:     req.CoverURL,
		Category:     req.Category,
		Tags:         req.Tags,
		Status:       status,
		RejectReason: "",
	}
	if err := logic.UpdatePost(post); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"id": strconv.FormatInt(id, 10)})
}

// PostDeleteHandler 删除文章；非管理员只能删除自己的文章。
func PostDeleteHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	userID := c.GetInt64(ContextUserIDKey)
	role := c.GetString(ContextUserRoleKey)
	post, err := logic.GetPostDetail(id)
	if err != nil {
		ResponseError(c, CodeNotFound)
		return
	}
	if role != models.RoleAdmin && post.AuthorID != userID {
		ResponseError(c, CodeAuthFailed)
		return
	}
	if err := logic.DeletePost(id); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"id": strconv.FormatInt(id, 10)})
}

// PostListHandler 文章列表（仅已审核通过），支持搜索与排序。
func PostListHandler(c *gin.Context) {
	sort := c.DefaultQuery("sort", "latest")
	keyword := c.DefaultQuery("keyword", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}
	items, total, err := logic.ListPosts(sort, keyword, page, size)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, PostListResponse{
		Items: items,
		Total: total,
	})
}

// PostTitlesHandler 文章标题列表，用于侧边栏或快速检索。
func PostTitlesHandler(c *gin.Context) {
	sort := c.DefaultQuery("sort", "latest")
	keyword := c.DefaultQuery("keyword", "")
	size, _ := strconv.Atoi(c.DefaultQuery("size", "50"))
	if size <= 0 || size > 200 {
		size = 50
	}
	items, err := logic.ListPostTitles(sort, keyword, size)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, items)
}

// PostDetailHandler 文章详情；未通过审核的文章仅作者或管理员可见。
func PostDetailHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	post, err := logic.GetPostDetail(id)
	if err != nil {
		ResponseError(c, CodeNotFound)
		return
	}
	if post.Status != models.StatusApproved {
		userID := c.GetInt64(ContextUserIDKey)
		role := c.GetString(ContextUserRoleKey)
		if role != models.RoleAdmin && post.AuthorID != userID {
			ResponseError(c, CodeNotFound)
			return
		}
	}
	ResponseSuccess(c, post)
}

// PostLikeHandler 点赞文章（仅已审核通过）。
func PostLikeHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	likeCount, err := logic.LikePost(id)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"likeCount": likeCount})
}

// MyPostListHandler 获取当前用户的文章列表。
func MyPostListHandler(c *gin.Context) {
	userID := c.GetInt64(ContextUserIDKey)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}
	items, total, err := logic.ListPostsByAuthor(userID, page, size)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, PostListResponse{
		Items: items,
		Total: total,
	})
}

// AdminPostListHandler 管理端按状态分页获取文章。
func AdminPostListHandler(c *gin.Context) {
	status := c.DefaultQuery("status", models.StatusPending)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}
	items, total, err := logic.ListPostsByStatus(status, page, size)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, PostListResponse{
		Items: items,
		Total: total,
	})
}

// AdminPostApproveHandler 管理员通过文章审核。
func AdminPostApproveHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	if err := logic.UpdatePostStatus(id, models.StatusApproved); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"id": strconv.FormatInt(id, 10), "status": models.StatusApproved})
}

// AdminPostRejectHandler 管理员驳回文章审核。
func AdminPostRejectHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	var req RejectRequest
	_ = c.ShouldBindJSON(&req)
	if err := logic.UpdatePostStatusWithReason(id, models.StatusRejected, req.Reason); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"id": strconv.FormatInt(id, 10), "status": models.StatusRejected})
}
