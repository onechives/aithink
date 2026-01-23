package mysql

import "aithink/models"

// CreateNicknameRequest 新建昵称申请记录。
func CreateNicknameRequest(req *models.NicknameRequest) error {
	query := "INSERT INTO nickname_requests (id, user_id, nickname, status, reject_reason, created_at, updated_at) VALUES (?, ?, ?, ?, '', NOW(), NOW())"
	_, err := db.Exec(query, req.ID, req.UserID, req.Nickname, req.Status)
	return err
}

// ListNicknameRequests 获取昵称申请列表（可按状态过滤）。
func ListNicknameRequests(status string) ([]models.NicknameRequest, error) {
	items := make([]models.NicknameRequest, 0)
	query := "SELECT id, user_id, nickname, status, reject_reason, created_at, updated_at FROM nickname_requests ORDER BY created_at DESC"
	args := []interface{}{}
	if status != "" {
		query = "SELECT id, user_id, nickname, status, reject_reason, created_at, updated_at FROM nickname_requests WHERE status=? ORDER BY created_at DESC"
		args = append(args, status)
	}
	if err := db.Select(&items, query, args...); err != nil {
		return nil, err
	}
	return items, nil
}

// UpdateNicknameRequestStatus 更新昵称申请状态（清空驳回原因）。
func UpdateNicknameRequestStatus(id int64, status string) error {
	_, err := db.Exec("UPDATE nickname_requests SET status=?, reject_reason='', updated_at=NOW() WHERE id=?", status, id)
	return err
}

// GetNicknameRequest 获取单条昵称申请。
func GetNicknameRequest(id int64) (*models.NicknameRequest, error) {
	item := new(models.NicknameRequest)
	query := "SELECT id, user_id, nickname, status, reject_reason, created_at, updated_at FROM nickname_requests WHERE id=?"
	if err := db.Get(item, query, id); err != nil {
		return nil, err
	}
	return item, nil
}

// UpdateNicknameRequestStatusWithReason 更新昵称申请状态并记录原因。
func UpdateNicknameRequestStatusWithReason(id int64, status, reason string) error {
	_, err := db.Exec("UPDATE nickname_requests SET status=?, reject_reason=?, updated_at=NOW() WHERE id=?", status, reason, id)
	return err
}

// IsNicknameRequested 判断昵称是否处于申请中。
func IsNicknameRequested(nickname string) (bool, error) {
	var count int64
	if err := db.Get(&count, "SELECT COUNT(*) FROM nickname_requests WHERE nickname=? AND status='pending'", nickname); err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasPendingNicknameRequest 判断用户是否已有待审核昵称申请。
func HasPendingNicknameRequest(userID int64) (bool, error) {
	var count int64
	if err := db.Get(&count, "SELECT COUNT(*) FROM nickname_requests WHERE user_id=? AND status='pending'", userID); err != nil {
		return false, err
	}
	return count > 0, nil
}
