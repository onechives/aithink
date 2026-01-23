package mysql

import (
	"aithink/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 不要忘了导入数据库驱动
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

// Init 初始化 MySQL 连接池。
func Init(cfg *settings.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName)

	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed, err", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(cfg.MaxOpenConnect) // 最大连接数
	db.SetMaxIdleConns(cfg.MaxIdleConnect) // 最大空闲连接数
	return
}

// Close 关闭数据库连接池。
func Close() {
	_ = db.Close()
}
