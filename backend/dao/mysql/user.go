package mysql

import (
	"aithink/models"
	"errors"
)

// CreateUser 新建用户（注册时使用）。
func CreateUser(user *models.User) error {
	query := "INSERT INTO users (id, username, nickname, password_hash, role, status, totp_secret, totp_enabled, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, '', 0, NOW(), NOW())"
	_, err := db.Exec(query, user.ID, user.Username, user.Nickname, user.Password, user.Role, user.Status)
	return err
}

// GetUserByUsername 按用户名查询用户。
func GetUserByUsername(username string) (*models.User, error) {
	user := new(models.User)
	query := "SELECT id, username, nickname, password_hash, role, status, reject_reason, totp_secret, totp_enabled, created_at, updated_at FROM users WHERE username=?"
	if err := db.Get(user, query, username); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID 按 ID 查询用户。
func GetUserByID(id int64) (*models.User, error) {
	user := new(models.User)
	query := "SELECT id, username, nickname, password_hash, role, status, reject_reason, totp_secret, totp_enabled, created_at, updated_at FROM users WHERE id=?"
	if err := db.Get(user, query, id); err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUserStatus 更新用户审核状态。
func UpdateUserStatus(id int64, status string) error {
	result, err := db.Exec("UPDATE users SET status=?, reject_reason='', updated_at=NOW() WHERE id=?", status, id)
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

// UpdateUserTOTP 更新用户 2FA 状态与密钥。
func UpdateUserTOTP(id int64, secret string, enabled bool) error {
	result, err := db.Exec("UPDATE users SET totp_secret=?, totp_enabled=?, updated_at=NOW() WHERE id=?", secret, enabled, id)
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

// ListUsersByStatus 按状态获取用户列表（管理端）。
func ListUsersByStatus(status string) ([]models.User, error) {
	items := make([]models.User, 0)
	query := "SELECT id, username, nickname, role, status, reject_reason, totp_enabled, created_at, updated_at FROM users WHERE status=? ORDER BY created_at DESC"
	if err := db.Select(&items, query, status); err != nil {
		return nil, err
	}
	return items, nil
}

// UpdateUserStatusWithReason 更新用户状态并记录拒绝原因。
func UpdateUserStatusWithReason(id int64, status, reason string) error {
	result, err := db.Exec("UPDATE users SET status=?, reject_reason=?, updated_at=NOW() WHERE id=?", status, reason, id)
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

// IsNicknameTaken 判断昵称是否已被占用。
func IsNicknameTaken(nickname string) (bool, error) {
	var count int64
	if err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE nickname=?", nickname); err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateUserNickname 更新用户昵称。
func UpdateUserNickname(id int64, nickname string) error {
	result, err := db.Exec("UPDATE users SET nickname=?, updated_at=NOW() WHERE id=?", nickname, id)
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
