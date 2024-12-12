package main

import (
	"context"
	"fmt"
	"github.com/allan7yin/rate-limiter/limiter"
	"github.com/redis/go-redis/v9"
	"time"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	bucketKey := "bitImageRateLimiter"
	var bucketMaxTokens int64 = 10
	var bucketRefillRate float64 = 1

	// reset Redis to consistent state before each run
	err := client.Del(ctx, bucketKey).Err()
	if err != nil {
		fmt.Println("Failed to clear Redis key:", err)
		return
	}

	rtb := limiter.NewRedisTokenBucket(client, bucketKey, bucketMaxTokens, bucketRefillRate)
	for i := 0; i < 15; i++ {
		allowed, err := rtb.AllowRequest(ctx, 1)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if allowed {
			fmt.Printf("Request %d: Allowed\n", i+1)
		} else {
			fmt.Printf("Request %d: Denied\n", i+1)
		}
		time.Sleep(300 * time.Millisecond)
	}
}
