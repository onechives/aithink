package middlewares

import (
	"aithink/controllers"
	"aithink/dao/redis"
	"aithink/pkg/jwt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWTAuthMiddleware 基于 JWT 的认证中间件：
// 1) 校验 Authorization Bearer
// 2) 校验 Redis 中的单点登录 token
// 3) 解析 JWT 并写入上下文
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controllers.ResponseError(c, controllers.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		// 单设备登录验证：判断请求 token 是否与 Redis 存储一致
		userId := strconv.FormatInt(mc.UserID, 10)
		redisToken, err := redis.GetTokenForRides(userId)
		if err != nil {
			zap.L().Error("没有在redis获取到这个token: ", zap.Error(err))
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		if redisToken != parts[1] {
			zap.L().Error("redis中的token和请求头中的token不一致: ", zap.Error(err))
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(controllers.ContextUserIDKey, mc.UserID)
		c.Set(controllers.ContextUserRoleKey, mc.Role)
		c.Next() // 后续的处理函数可以用过c.Get("userID")来获取当前请求的用户信息
	}
}

// OptionalJWTAuthMiddleware 如果请求携带 token，则校验并写入上下文；
// 用于既可匿名访问又可识别登录态的接口（如文章详情）。
func OptionalJWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}

		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		userId := strconv.FormatInt(mc.UserID, 10)
		redisToken, err := redis.GetTokenForRides(userId)
		if err != nil || redisToken != parts[1] {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		c.Set(controllers.ContextUserIDKey, mc.UserID)
		c.Set(controllers.ContextUserRoleKey, mc.Role)
		c.Next()
	}
}

// AdminOnlyMiddleware 仅管理员可访问
func AdminOnlyMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		role, ok := c.Get(controllers.ContextUserRoleKey)
		if !ok || role != "admin" {
			controllers.ResponseError(c, controllers.CodeAuthFailed)
			c.Abort()
			return
		}
		c.Next()
	}
}
