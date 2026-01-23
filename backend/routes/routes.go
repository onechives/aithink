package routes

import (
	"aithink/controllers"
	"aithink/logger"
	"aithink/middlewares"
	"aithink/settings"
	"time"

	"github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"
)

// SetUp 初始化路由、中间件与分组；返回可直接启动的 Gin 引擎实例。
func SetUp(cfg *settings.AppConfig) *gin.Engine {

	r := gin.New()

	// 1、如果有 Nginx / LB /docker
	// r.SetTrustedProxies([]string{
	// 	"127.0.0.1",     //nginx
	// 	// "10.0.0.0/8",    //lb
	// 	// "172.16.0.0/12", //docker
	// })

	//nginx里面还需要配置：
	// proxy_set_header X-Real-IP $remote_addr;
	// proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

	// 2) 未部署 Nginx 时，关闭代理信任（避免伪造来源）
	r.SetTrustedProxies(nil)
	r.MaxMultipartMemory = 8 << 20

	limiter := middlewares.NewLimiter()
	// 日志、异常恢复、跨域与限流中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.CORS(), limiter.Middleware(10*time.Second, 1000)) //每2秒添加10个令牌
	r.Static("/uploads", "./uploads")

	v1 := r.Group("/api/v1")
	// 公开接口：注册/登录/文章列表与详情
	v1.POST("/register", controllers.RegisterHandler)
	v1.POST("/login", controllers.LoginHandler)
	v1.POST("/login/verify", controllers.LoginVerifyHandler)

	v1.GET("/posts", controllers.PostListHandler)
	v1.GET("/post-titles", controllers.PostTitlesHandler)
	v1.GET("/posts/:id", middlewares.OptionalJWTAuthMiddleware(), controllers.PostDetailHandler)
	v1.POST("/posts/:id/like", controllers.PostLikeHandler)

	// 需登录访问的接口
	auth := v1.Group("")
	auth.Use(middlewares.JWTAuthMiddleware())
	{
		auth.GET("/me", controllers.MeHandler)
		auth.GET("/me/posts", controllers.MyPostListHandler)
		auth.GET("/me/messages", controllers.MessageListHandler)
		auth.GET("/me/messages/unread-count", controllers.MessageUnreadCountHandler)
		auth.POST("/me/messages/:id/read", controllers.MessageReadHandler)
		auth.POST("/me/2fa/init", controllers.TOTPInitHandler)
		auth.POST("/me/2fa/enable", controllers.TOTPEnableHandler)
		auth.POST("/me/2fa/disable", controllers.TOTPDisableHandler)
		auth.POST("/me/nickname", controllers.NicknameRequestHandler)

		auth.POST("/upload", controllers.UploadImageHandler)
		auth.POST("/posts", controllers.PostCreateHandler)
		auth.PUT("/posts/:id", controllers.PostUpdateHandler)
		auth.DELETE("/posts/:id", controllers.PostDeleteHandler)
	}

	// 管理员接口：用户/文章/昵称审核
	admin := v1.Group("/admin")
	admin.Use(middlewares.JWTAuthMiddleware(), middlewares.AdminOnlyMiddleware())
	{
		admin.GET("/users", controllers.AdminUserListHandler)
		admin.POST("/users/:id/approve", controllers.AdminUserApproveHandler)
		admin.POST("/users/:id/reject", controllers.AdminUserRejectHandler)

		admin.GET("/posts", controllers.AdminPostListHandler)
		admin.POST("/posts/:id/approve", controllers.AdminPostApproveHandler)
		admin.POST("/posts/:id/reject", controllers.AdminPostRejectHandler)

		admin.GET("/nicknames", controllers.AdminNicknameListHandler)
		admin.POST("/nicknames/:id/approve", controllers.AdminNicknameApproveHandler)
		admin.POST("/nicknames/:id/reject", controllers.AdminNicknameRejectHandler)
	}

	pprof.Register(r) //Go 官方性能分析工具 http://localhost:8080/debug/pprof/
	return r
}
