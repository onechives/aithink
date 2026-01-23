package logic

import (
	"aithink/dao/mysql"
	"aithink/models"
)

// ListPendingUsers 获取待审核用户列表（旧接口保留）。
func ListPendingUsers() ([]models.User, error) {
	return mysql.ListUsersByStatus(models.StatusPending)
}

// ListUsersByStatus 按状态获取用户列表。
func ListUsersByStatus(status string) ([]models.User, error) {
	return mysql.ListUsersByStatus(status)
}

// ApproveUser 通过用户审核，并发送站内信通知。
func ApproveUser(userID int64) error {
	if err := mysql.UpdateUserStatus(userID, models.StatusApproved); err != nil {
		return err
	}
	return SendMessage(userID, "注册审核通过", "你的账号已通过审核，可以登录并发布文章。")
}

// RejectUser 驳回用户审核，并发送站内信通知。
func RejectUser(userID int64, reason string) error {
	if err := mysql.UpdateUserStatusWithReason(userID, models.StatusRejected, reason); err != nil {
		return err
	}
	content := "你的账号审核未通过，请联系管理员。"
	if reason != "" {
		content = content + " 原因：" + reason
	}
	return SendMessage(userID, "注册审核未通过", content)
}
