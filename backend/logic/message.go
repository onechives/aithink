package logic

import (
	"aithink/dao/mysql"
	"aithink/models"
	"aithink/pkg/snowflake"
)

// SendMessage 发送站内信（用于审核结果与系统通知）。
func SendMessage(userID int64, title, content string) error {
	msg := &models.Message{
		ID:      snowflake.GenIDByInt(),
		UserID:  userID,
		Title:   title,
		Content: content,
		Status:  models.MessageUnread,
	}
	return mysql.CreateMessage(msg)
}
