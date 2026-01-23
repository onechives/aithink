package mysql

import "aithink/models"

// CreateMessage 新增站内信记录。
func CreateMessage(msg *models.Message) error {
	query := "INSERT INTO messages (id, user_id, title, content, status, created_at) VALUES (?, ?, ?, ?, ?, NOW())"
	_, err := db.Exec(query, msg.ID, msg.UserID, msg.Title, msg.Content, msg.Status)
	return err
}

// ListMessages 获取用户消息列表（可按状态过滤）。
func ListMessages(userID int64, status string) ([]models.Message, error) {
	items := make([]models.Message, 0)
	query := "SELECT id, user_id, title, content, status, created_at FROM messages WHERE user_id=? ORDER BY created_at DESC"
	args := []interface{}{userID}
	if status != "" {
		query = "SELECT id, user_id, title, content, status, created_at FROM messages WHERE user_id=? AND status=? ORDER BY created_at DESC"
		args = []interface{}{userID, status}
	}
	if err := db.Select(&items, query, args...); err != nil {
		return nil, err
	}
	return items, nil
}

// MarkMessageRead 标记消息为已读。
func MarkMessageRead(id int64, userID int64) error {
	_, err := db.Exec("UPDATE messages SET status='read' WHERE id=? AND user_id=?", id, userID)
	return err
}

// CountUnreadMessages 统计未读消息数量。
func CountUnreadMessages(userID int64) (int64, error) {
	var count int64
	if err := db.Get(&count, "SELECT COUNT(*) FROM messages WHERE user_id=? AND status='unread'", userID); err != nil {
		return 0, err
	}
	return count, nil
}
