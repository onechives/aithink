package main

import (
	"aithink/controllers"
	"aithink/dao/mysql"
	"aithink/dao/redis"
	"aithink/logger"
	"aithink/pkg/snowflake"
	"aithink/routes"
	"aithink/settings"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 1) 加载配置：支持命令行传参指定配置文件路径
	var filename = "./conf/config.yaml" //没有带配置文件参数 就用默认的
	//命令行参数获取配置文件路径
	if len(os.Args) > 1 { //如果有带配置文件参数就用参数的配置
		filename = os.Args[1] //第一个参数
	}
	if err := settings.Init(filename); err != nil {
		fmt.Println("配置文件加载失败！", err.Error())
		return
	} else {
		fmt.Println("配置文件加载成功！")
	}

	// 2) 初始化日志：统一结构化输出，便于排查问题
	if err := logger.Init(settings.Config.LogConfig, settings.Config.Mode); err != nil {
		fmt.Println("日志系统初始化失败！", err)
		return
	} else {
		fmt.Println("日志系统初始化成功！")
	}
	defer zap.L().Sync()

	// 3) 初始化 MySQL 连接池
	if err := mysql.Init(settings.Config.MysqlConfig); err != nil {
		zap.L().Error("mysql初始化失败！", zap.Error(err))
		return
	} else {
		zap.L().Info("mysql初始化成功！")
	}
	defer mysql.Close()

	// 4) 初始化 Redis 连接池，用于单点登录与缓存
	if err := redis.Init(settings.Config.RedisConfig); err != nil {
		zap.L().Error("redis初始化失败！", zap.Error(err))
		return
	} else {
		zap.L().Info("redis初始化成功！")
	}
	defer redis.Close()

	// 5) 初始化雪花算法：生成分布式唯一 ID
	if err := snowflake.Init(settings.Config.StartTime, settings.Config.MachineID); err != nil {
		fmt.Printf("init failed, err:%v\n", err)
		panic(err)
	}
	snowId := snowflake.GenIDByString()
	fmt.Println("snowId=", snowId)

	// 6) 初始化 gin 校验器的翻译器（中文提示）
	if err := controllers.InitTrans("zh"); err != nil {
		zap.L().Error(err.Error())
		return
	}

	// 7) 注册路由与中间件
	r := routes.SetUp(settings.Config)

	// 8) 启动 HTTP 服务，并支持优雅关机
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Config.Port),
		Handler: r,
	}
	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("listen error:", zap.Error(err))
		}
	}()
	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown，服务关机: ", zap.Error(err))
	}
	zap.L().Info("Server exiting，服务退出")
}
