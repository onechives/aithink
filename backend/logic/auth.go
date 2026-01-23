package logic

import (
	"aithink/dao/mysql"
	"aithink/dao/redis"
	"aithink/models"
	"aithink/pkg/jwt"
	"aithink/pkg/snowflake"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var (
	// 认证/账号相关错误
	ErrAuthFailed      = errors.New("auth failed")
	ErrUserPending     = errors.New("user pending approval")
	ErrUserRejected    = errors.New("user rejected")
	ErrInvalidTwoFA    = errors.New("invalid 2fa")
	ErrTwoFAInitFirst  = errors.New("2fa not initialized")
	ErrNicknameInvalid = errors.New("nickname invalid")
	ErrNicknameTaken   = errors.New("nickname taken")
)

// tempLoginTTL 登录二次验证临时 token 有效期
const tempLoginTTL = 5 * time.Minute
// totpSetupTTL 2FA 初始化密钥在 Redis 的有效期
const totpSetupTTL = 10 * time.Minute

// RegisterUser 注册新用户，默认待审核状态。
func RegisterUser(username, password string) error {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return ErrAuthFailed
	}
	if ok, err := mysql.IsNicknameTaken(username); err != nil {
		return err
	} else if ok {
		return ErrAuthFailed
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &models.User{
		ID:       snowflake.GenIDByInt(),
		Username: username,
		Nickname: username,
		Password: string(hashed),
		Role:     models.RoleUser,
		Status:   models.StatusPending,
	}
	return mysql.CreateUser(user)
}

// LoginStep 登录第一步：校验账号密码；若开启 2FA 则发临时 token。
func LoginStep(username, password string) (userID int64, role string, token string, need2FA bool, tempToken string, err error) {
	user, err := mysql.GetUserByUsername(username)
	if err != nil {
		return 0, "", "", false, "", ErrAuthFailed
	}
	if user.Status == models.StatusPending {
		return 0, "", "", false, "", ErrUserPending
	}
	if user.Status == models.StatusRejected {
		return 0, "", "", false, "", ErrUserRejected
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, "", "", false, "", ErrAuthFailed
	}

	if user.TOTPEnabled {
		tt, err := generateTempToken()
		if err != nil {
			return 0, "", "", false, "", err
		}
		if err := redis.SetTempLoginToken(tt, formatUserID(user.ID), tempLoginTTL); err != nil {
			return 0, "", "", false, "", err
		}
		return user.ID, user.Role, "", true, tt, nil
	}

	token, err = issueToken(user)
	if err != nil {
		return 0, "", "", false, "", err
	}
	return user.ID, user.Role, token, false, "", nil
}

// LoginVerify 登录第二步：验证 2FA，成功后签发正式 token。
func LoginVerify(tempToken, code string) (userID int64, role string, token string, err error) {
	userIDStr, err := redis.GetTempLoginUserID(tempToken)
	if err != nil {
		return 0, "", "", ErrAuthFailed
	}
	userID, err = parseUserID(userIDStr)
	if err != nil {
		return 0, "", "", ErrAuthFailed
	}
	user, err := mysql.GetUserByID(userID)
	if err != nil {
		return 0, "", "", ErrAuthFailed
	}
	if !user.TOTPEnabled {
		return 0, "", "", ErrAuthFailed
	}
	if ok := totp.Validate(code, user.TOTPSecret); !ok {
		return 0, "", "", ErrInvalidTwoFA
	}
	_ = redis.DeleteTempLoginToken(tempToken)
	token, err = issueToken(user)
	if err != nil {
		return 0, "", "", err
	}
	return user.ID, user.Role, token, nil
}

// IssueUserToken 根据用户 ID 直接签发 token（用于系统内刷新或特殊场景）。
func IssueUserToken(userID int64) (string, error) {
	user, err := mysql.GetUserByID(userID)
	if err != nil {
		return "", err
	}
	return issueToken(user)
}

// InitTOTP 初始化 2FA，生成密钥与 otpauth URL。
func InitTOTP(userID int64, username string) (secret string, otpauth string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "SimpleBlog",
		AccountName: username,
	})
	if err != nil {
		return "", "", err
	}
	if err := redis.SetTOTPSetupSecret(formatUserID(userID), key.Secret(), totpSetupTTL); err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

// EnableTOTP 验证 2FA 密码后启用 TOTP。
func EnableTOTP(userID int64, code string) error {
	secret, err := redis.GetTOTPSetupSecret(formatUserID(userID))
	if err != nil {
		return ErrTwoFAInitFirst
	}
	if ok := totp.Validate(code, secret); !ok {
		return ErrInvalidTwoFA
	}
	if err := mysql.UpdateUserTOTP(userID, secret, true); err != nil {
		return err
	}
	_ = redis.DeleteTOTPSetupSecret(formatUserID(userID))
	return nil
}

// DisableTOTP 验证 2FA 密码后关闭 TOTP。
func DisableTOTP(userID int64, code string) error {
	user, err := mysql.GetUserByID(userID)
	if err != nil {
		return err
	}
	if !user.TOTPEnabled {
		return nil
	}
	if ok := totp.Validate(code, user.TOTPSecret); !ok {
		return ErrInvalidTwoFA
	}
	return mysql.UpdateUserTOTP(userID, "", false)
}

// GetUserByID 获取用户信息（供控制器使用）。
func GetUserByID(userID int64) (*models.User, error) {
	return mysql.GetUserByID(userID)
}

// RequestNicknameChange 申请昵称变更（进入审核队列）。
func RequestNicknameChange(userID int64, nickname string) error {
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return ErrNicknameInvalid
	}
	if utf8.RuneCountInString(nickname) > 10 {
		return ErrNicknameInvalid
	}
	if ok, err := mysql.HasPendingNicknameRequest(userID); err != nil {
		return err
	} else if ok {
		return ErrNicknameInvalid
	}
	if ok, err := mysql.IsNicknameTaken(nickname); err != nil {
		return err
	} else if ok {
		return ErrNicknameTaken
	}
	if ok, err := mysql.IsNicknameRequested(nickname); err != nil {
		return err
	} else if ok {
		return ErrNicknameTaken
	}
	req := &models.NicknameRequest{
		ID:       snowflake.GenIDByInt(),
		UserID:   userID,
		Nickname: nickname,
		Status:   models.NicknamePending,
	}
	return mysql.CreateNicknameRequest(req)
}

// ApproveNicknameRequest 通过昵称审核并写入用户表，同时发送通知。
func ApproveNicknameRequest(id int64) error {
	req, err := mysql.GetNicknameRequest(id)
	if err != nil {
		return err
	}
	if err := mysql.UpdateNicknameRequestStatus(id, models.NicknameApproved); err != nil {
		return err
	}
	if err := mysql.UpdateUserNickname(req.UserID, req.Nickname); err != nil {
		return err
	}
	return SendMessage(req.UserID, "昵称审核通过", "你的昵称已更新为："+req.Nickname)
}

// RejectNicknameRequest 驳回昵称审核并记录原因，同时发送通知。
func RejectNicknameRequest(id int64, reason string) error {
	req, err := mysql.GetNicknameRequest(id)
	if err != nil {
		return err
	}
	if err := mysql.UpdateNicknameRequestStatusWithReason(id, models.NicknameRejected, reason); err != nil {
		return err
	}
	content := "你的昵称审核未通过，请修改后重新提交。"
	if reason != "" {
		content = content + " 原因：" + reason
	}
	return SendMessage(req.UserID, "昵称审核未通过", content)
}

// ListNicknameRequests 获取昵称申请列表（支持状态筛选）。
func ListNicknameRequests(status string) ([]models.NicknameRequest, error) {
	return mysql.ListNicknameRequests(status)
}

// generateTempToken 生成短期登录临时 token。
func generateTempToken() (string, error) {
	raw := make([]byte, 24)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

// issueToken 签发 JWT 并写入 Redis，用于单点登录校验。
func issueToken(user *models.User) (string, error) {
	token, err := jwt.GenToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", err
	}
	if err := redis.SetTokenForUser(formatUserID(user.ID), token, jwt.TokenExpireDuration); err != nil {
		return "", err
	}
	return token, nil
}

// formatUserID 将用户 ID 统一转换为字符串（Redis key 使用）。
func formatUserID(id int64) string {
	return strconv.FormatInt(id, 10)
}

// parseUserID 解析字符串用户 ID。
func parseUserID(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
