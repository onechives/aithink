package controllers

import (
	"aithink/dao/mysql"
	"aithink/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type MessageListResponse struct {
	// 消息列表包装
	Items []MessageItemResponse `json:"items"`
}

type MessageItemResponse struct {
	// 消息基础字段 + 可选的国际化标记
	ID          int64             `json:"id"`
	UserID      int64             `json:"userId"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"createdAt"`
	MessageType string            `json:"messageType,omitempty"`
	I18nKey     string            `json:"i18nKey,omitempty"`
	Params      map[string]string `json:"params,omitempty"`
}

// MessageListHandler 获取当前用户的站内信列表，支持状态过滤。
func MessageListHandler(c *gin.Context) {
	userID := c.GetInt64(ContextUserIDKey)
	if userID == 0 {
		ResponseError(c, CodeNeedLogin)
		return
	}
	status := c.DefaultQuery("status", "")
	items, err := mysql.ListMessages(userID, status)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	respItems := make([]MessageItemResponse, 0, len(items))
	for _, item := range items {
		msgType, i18nKey, params := mapMessageI18n(item)
		respItems = append(respItems, MessageItemResponse{
			ID:          item.ID,
			UserID:      item.UserID,
			Title:       item.Title,
			Content:     item.Content,
			Status:      item.Status,
			CreatedAt:   item.CreatedAt,
			MessageType: msgType,
			I18nKey:     i18nKey,
			Params:      params,
		})
	}
	ResponseSuccess(c, MessageListResponse{Items: respItems})
}

// MessageReadHandler 标记单条消息为已读。
func MessageReadHandler(c *gin.Context) {
	userID := c.GetInt64(ContextUserIDKey)
	if userID == 0 {
		ResponseError(c, CodeNeedLogin)
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	if err := mysql.MarkMessageRead(id, userID); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"id": id, "status": models.MessageRead})
}

// MessageUnreadCountHandler 获取未读消息数量。
func MessageUnreadCountHandler(c *gin.Context) {
	userID := c.GetInt64(ContextUserIDKey)
	if userID == 0 {
		ResponseError(c, CodeNeedLogin)
		return
	}
	count, err := mysql.CountUnreadMessages(userID)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{"count": count})
}

// mapMessageI18n 将后端消息标题映射为前端 i18n key，便于多语言渲染。
func mapMessageI18n(item models.Message) (string, string, map[string]string) {
	params := map[string]string{}
	switch item.Title {
	case "注册审核通过":
		return "user_review", "messages.userApproved", nil
	case "注册审核未通过":
		if reason := extractReason(item.Content); reason != "" {
			params["reason"] = reason
		}
		return "user_review", "messages.userRejected", params
	case "昵称审核通过":
		if nickname := strings.TrimPrefix(item.Content, "你的昵称已更新为："); nickname != item.Content {
			params["nickname"] = nickname
		}
		return "nickname_review", "messages.nicknameApproved", params
	case "昵称审核未通过":
		if reason := extractReason(item.Content); reason != "" {
			params["reason"] = reason
		}
		return "nickname_review", "messages.nicknameRejected", params
	case "文章审核通过":
		return "post_review", "messages.postApproved", nil
	case "文章审核未通过":
		if reason := extractReason(item.Content); reason != "" {
			params["reason"] = reason
		}
		return "post_review", "messages.postRejected", params
	default:
		return "", "", nil
	}
}

// extractReason 从消息内容中抽取“原因”字段，便于国际化插值。
func extractReason(content string) string {
	parts := strings.SplitN(content, " 原因：", 2)
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
