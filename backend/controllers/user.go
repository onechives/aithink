package controllers

import (
	"aithink/logic"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	// 登录账号（用户名）
	Username string `json:"username" binding:"required"`
	// 登录密码（明文入参，后端自行校验）
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	// 用户 ID
	UserID    int64  `json:"userId"`
	// 用户角色：admin / user
	Role      string `json:"role"`
	// 登录成功直接下发的 token（若开启 2FA 则为空）
	Token     string `json:"token,omitempty"`
	// 是否需要两步验证
	Need2FA   bool   `json:"need2fa"`
	// 两步验证临时 token，用于二次校验
	TempToken string `json:"tempToken,omitempty"`
}

type RegisterRequest struct {
	// 注册用户名
	Username string `json:"username" binding:"required"`
	// 注册密码
	Password string `json:"password" binding:"required"`
}

type LoginVerifyRequest struct {
	// 登录时下发的临时 token
	TempToken string `json:"tempToken" binding:"required"`
	// TOTP 验证码
	Code      string `json:"code" binding:"required"`
}

type MeResponse struct {
	// 用户基础信息与安全状态
	UserID      int64  `json:"userId"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	TOTPEnabled bool   `json:"totpEnabled"`
}

type TOTPInitResponse struct {
	// 2FA 密钥与 OTP URL（前端生成二维码）
	Secret string `json:"secret"`
	URL    string `json:"url"`
}

type TOTPVerifyRequest struct {
	// 2FA 校验码
	Code string `json:"code" binding:"required"`
}

// RegisterHandler 处理用户注册：写入待审核账号。
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	if err := logic.RegisterUser(req.Username, req.Password); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"status": "pending"})
}

// LoginHandler 登录第一步：校验账号密码，若开启 2FA 则下发临时 token。
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	userID, role, token, need2FA, tempToken, err := logic.LoginStep(req.Username, req.Password)
	if err != nil {
		if err == logic.ErrUserPending {
			ResponseErrorWithMsg(c, CodeUserPending, "账号待审核")
			return
		}
		if err == logic.ErrUserRejected {
			ResponseErrorWithMsg(c, CodeAuthFailed, "账号审核未通过")
			return
		}
		ResponseError(c, CodeAuthFailed)
		return
	}
	c.JSON(http.StatusOK, ResponseData{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data: LoginResponse{
			UserID:    userID,
			Role:      role,
			Token:     token,
			Need2FA:   need2FA,
			TempToken: tempToken,
		},
	})
}

// LoginVerifyHandler 登录第二步：验证 TOTP，成功后返回正式 token。
func LoginVerifyHandler(c *gin.Context) {
	var req LoginVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	userID, role, token, err := logic.LoginVerify(req.TempToken, req.Code)
	if err != nil {
		ResponseErrorWithMsg(c, CodeAuthFailed, "两步验证失败")
		return
	}
	c.JSON(http.StatusOK, ResponseData{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data: LoginResponse{
			UserID:  userID,
			Role:    role,
			Token:   token,
			Need2FA: false,
		},
	})
}

// MeHandler 获取当前登录用户的基础信息与安全设置。
func MeHandler(c *gin.Context) {
	userID := c.GetInt64(ContextUserIDKey)
	if userID == 0 {
		ResponseError(c, CodeNeedLogin)
		return
	}
	user, err := logic.GetUserByID(userID)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, MeResponse{
		UserID:      user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Role:        user.Role,
		Status:      user.Status,
		TOTPEnabled: user.TOTPEnabled,
	})
}

// TOTPInitHandler 初始化 2FA：生成密钥与二维码地址。
func TOTPInitHandler(c *gin.Context) {
	userID := c.GetInt64(ContextUserIDKey)
	if userID == 0 {
		ResponseError(c, CodeNeedLogin)
		return
	}
	user, err := logic.GetUserByID(userID)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	secret, url, err := logic.InitTOTP(userID, user.Username)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, TOTPInitResponse{Secret: secret, URL: url})
}

// TOTPEnableHandler 启用 2FA：校验验证码并写入启用状态。
func TOTPEnableHandler(c *gin.Context) {
	userID := c.GetInt64(ContextUserIDKey)
	if userID == 0 {
		ResponseError(c, CodeNeedLogin)
		return
	}
	var req TOTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	if err := logic.EnableTOTP(userID, req.Code); err != nil {
		ResponseErrorWithMsg(c, CodeAuthFailed, "两步验证码错误")
		return
	}
	ResponseSuccess(c, gin.H{"enabled": true})
}

// TOTPDisableHandler 关闭 2FA：校验验证码并清除启用状态。
func TOTPDisableHandler(c *gin.Context) {
	userID := c.GetInt64(ContextUserIDKey)
	if userID == 0 {
		ResponseError(c, CodeNeedLogin)
		return
	}
	var req TOTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	if err := logic.DisableTOTP(userID, req.Code); err != nil {
		ResponseErrorWithMsg(c, CodeAuthFailed, "两步验证码错误")
		return
	}
	ResponseSuccess(c, gin.H{"enabled": false})
}

type NicknameRequest struct {
	// 新昵称（需审核）
	Nickname string `json:"nickname" binding:"required"`
}

// NicknameRequestHandler 申请昵称变更（进入审核流程）。
func NicknameRequestHandler(c *gin.Context) {
	userID := c.GetInt64(ContextUserIDKey)
	if userID == 0 {
		ResponseError(c, CodeNeedLogin)
		return
	}
	var req NicknameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	if err := logic.RequestNicknameChange(userID, req.Nickname); err != nil {
		if err == logic.ErrNicknameTaken {
			ResponseError(c, CodeNicknameTaken)
			return
		}
		ResponseError(c, CodeInvalidParams)
		return
	}
	ResponseSuccess(c, gin.H{"status": "pending"})
}
