package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// // RateLimitMiddleware 令牌桶限流中间件
// func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
// 	bucket := ratelimit.NewBucket(fillInterval, cap)
// 	return func(c *gin.Context) {
// 		if bucket.TakeAvailable(1) == 0 {
// 			c.String(http.StatusTooManyRequests, "太多请求了 限流中...Too many requests...")
// 			c.Abort()
// 			return
// 		}
// 		c.Next() //取到令牌就放行
// 	}
// }

type Limiter struct {
	// buckets 以客户端 IP 为维度存储限流桶
	buckets sync.Map // map[string]*ratelimit.Bucket
}

// NewLimiter 创建一个可复用的限流器实例。
func NewLimiter() *Limiter {
	return &Limiter{}
}

// Middleware 为每个客户端 IP 分配令牌桶，超过配额直接返回 429。
func (l *Limiter) Middleware(fillInterval time.Duration, cap int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 以 IP 为 key，确保单个 IP 限流互不影响
		key := c.ClientIP()
		bucketIface, _ := l.buckets.LoadOrStore(
			key,
			ratelimit.NewBucket(fillInterval, cap),
		)
		bucket := bucketIface.(*ratelimit.Bucket)

		if bucket.TakeAvailable(1) == 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求太多...too many requests...",
			})
			return
		}
		c.Next()
	}
}
