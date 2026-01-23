package controllers

import (
	"aithink/logic"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminNicknameListHandler 管理端获取昵称申请列表。
func AdminNicknameListHandler(c *gin.Context) {
	status := c.DefaultQuery("status", "")
	items, err := logic.ListNicknameRequests(status)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"items": items})
}

// AdminNicknameApproveHandler 通过昵称审核并写入用户表。
func AdminNicknameApproveHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	if err := logic.ApproveNicknameRequest(id); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"id": id, "status": "approved"})
}

// AdminNicknameRejectHandler 驳回昵称审核并记录原因。
func AdminNicknameRejectHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	var req RejectRequest
	_ = c.ShouldBindJSON(&req)
	if err := logic.RejectNicknameRequest(id, req.Reason); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"id": id, "status": "rejected"})
}
