package main

import (
	"context"
	"fmt"
	"github.com/allan7yin/rate-limiter/config"
	"github.com/allan7yin/rate-limiter/limiter"
	"github.com/allan7yin/rate-limiter/server"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

func RateLimiterMiddleware(rtb *limiter.RedisTokenBucket) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			allowed, err := rtb.AllowRequest(ctx, 1)
			if err != nil {
				log.Printf("Error checking rate limiter: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				log.Println("Too Many Requests")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	ctx := context.Background()
	c := config.LoadConfig()
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%s", c.RedisPort),
	})

	bucketKey := c.BucketKey
	bucketMaxTokens := c.BucketMaxTokens
	bucketRefillRate := c.BucketRefillRate

	err := client.Del(ctx, bucketKey).Err()
	if err != nil {
		fmt.Println("Failed to clear Redis key:", err)
		return
	}

	rtb := limiter.NewRedisTokenBucket(client, bucketKey, bucketMaxTokens, bucketRefillRate)

	s := server.NewServer()
	port := "localhost:" + c.AppPort

	rateLimitedHandler := RateLimiterMiddleware(rtb)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Request allowed!"))
	}))

	// TODO: Add handler to make API call to ImageStore
	s.AddRoute("/v1/test", rateLimitedHandler)

	err = s.Start(port)
	if err != nil {
		return
	}
}
