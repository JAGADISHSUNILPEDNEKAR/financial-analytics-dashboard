package middleware

import (
    "fmt"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "golang.org/x/time/rate"
)

type RateLimiter interface {
    Allow(key string) bool
}

type RedisRateLimiter struct {
    client *redis.Client
    rate   int
    window time.Duration
}

func NewRateLimiter(client *redis.Client) RateLimiter {
    return &RedisRateLimiter{
        client: client,
        rate:   100, // 100 requests
        window: time.Minute,
    }
}

func (r *RedisRateLimiter) Allow(key string) bool {
    ctx := context.Background()
    now := time.Now()
    windowStart := now.Add(-r.window)
    
    pipe := r.client.Pipeline()
    
    // Remove old entries
    pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))
    
    // Count current window
    count := pipe.ZCard(ctx, key)
    
    // Add current request
    pipe.ZAdd(ctx, key, redis.Z{
        Score:  float64(now.UnixNano()),
        Member: now.UnixNano(),
    })
    
    // Set expiration
    pipe.Expire(ctx, key, r.window)
    
    _, err := pipe.Exec(ctx)
    if err != nil {
        return false
    }
    
    return count.Val() < int64(r.rate)
}

func RateLimit(limiter RateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        key := fmt.Sprintf("rate_limit:%s:%s", userID, c.Request.URL.Path)
        
        if !limiter.Allow(key) {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": 60,
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}