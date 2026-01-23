package redis

import (
	"aithink/settings"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var rdb *redis.Client

// Redis key 前缀约定
const tokenKeyPrefix = "blog:token:"
const loginTempKeyPrefix = "blog:login:temp:"
const totpSetupKeyPrefix = "blog:totp:setup:"

// Init 初始化 Redis 客户端。
func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), //地址端口
		Password: cfg.Password,                             // 密码
		DB:       cfg.DbName,                               // 数据库
		PoolSize: cfg.PoolSize,                             // 连接池大小
	})

	_, err = rdb.Ping().Result()
	return err
}

func Close() {
	_ = rdb.Close()
}

// SetTokenForUser 保存用户登录 token（用于单点登录校验）。
func SetTokenForUser(userID string, token string, ttl time.Duration) error {
	return rdb.Set(tokenKeyPrefix+userID, token, ttl).Err()
}

// GetTokenForRides 获取用户当前登录 token。
func GetTokenForRides(userID string) (string, error) {
	return rdb.Get(tokenKeyPrefix + userID).Result()
}

// DeleteTokenForUser 删除用户登录 token。
func DeleteTokenForUser(userID string) error {
	return rdb.Del(tokenKeyPrefix + userID).Err()
}

// SetTempLoginToken 保存 2FA 登录临时 token。
func SetTempLoginToken(token string, userID string, ttl time.Duration) error {
	return rdb.Set(loginTempKeyPrefix+token, userID, ttl).Err()
}

// GetTempLoginUserID 获取临时 token 对应的用户 ID。
func GetTempLoginUserID(token string) (string, error) {
	return rdb.Get(loginTempKeyPrefix + token).Result()
}

// DeleteTempLoginToken 删除临时 token。
func DeleteTempLoginToken(token string) error {
	return rdb.Del(loginTempKeyPrefix + token).Err()
}

// SetTOTPSetupSecret 保存 2FA 初始化密钥（短期有效）。
func SetTOTPSetupSecret(userID string, secret string, ttl time.Duration) error {
	return rdb.Set(totpSetupKeyPrefix+userID, secret, ttl).Err()
}

// GetTOTPSetupSecret 获取 2FA 初始化密钥。
func GetTOTPSetupSecret(userID string) (string, error) {
	return rdb.Get(totpSetupKeyPrefix + userID).Result()
}

// DeleteTOTPSetupSecret 删除 2FA 初始化密钥。
func DeleteTOTPSetupSecret(userID string) error {
	return rdb.Del(totpSetupKeyPrefix + userID).Err()
}
