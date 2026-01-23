package controllers

import (
	"aithink/logic"
	"aithink/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	// 管理端用户列表
	Items []models.User `json:"items"`
}

type RejectRequest struct {
	// 审核驳回原因
	Reason string `json:"reason"`
}

// AdminUserListHandler 管理端获取用户列表（按状态过滤）。
func AdminUserListHandler(c *gin.Context) {
	status := c.DefaultQuery("status", models.StatusPending)
	items, err := logic.ListUsersByStatus(status)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, UserListResponse{Items: items})
}

// AdminUserApproveHandler 管理端通过用户注册审核。
func AdminUserApproveHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	if err := logic.ApproveUser(id); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"id": id, "status": models.StatusApproved})
}

// AdminUserRejectHandler 管理端驳回用户注册审核。
func AdminUserRejectHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	var req RejectRequest
	_ = c.ShouldBindJSON(&req)
	if err := logic.RejectUser(id, req.Reason); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"id": id, "status": models.StatusRejected})
}
