package limiter

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisTokenBucket struct {
	client     *redis.Client
	maxTokens  int64
	refillRate float64
	bucketKey  string
}

func NewRedisTokenBucket(client *redis.Client, key string, maxTokens int64, refillRate float64) *RedisTokenBucket {
	return &RedisTokenBucket{
		client:     client,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		bucketKey:  key,
	}
}

func (rtb *RedisTokenBucket) AllowRequest(ctx context.Context, tokens float64) (bool, error) {
	script := `
			local bucket = redis.call("HGETALL", KEYS[1])
			local maxTokens = tonumber(ARGV[1])
			local refillRate = tonumber(ARGV[2])
			local tokensRequested = tonumber(ARGV[3])
			local now = tonumber(ARGV[4])
			
			local currentTokens = maxTokens
			local lastRefillTime = now
			
			if next(bucket) ~= nil then
				lastRefillTime = tonumber(bucket[2])
				currentTokens = tonumber(bucket[4])
			
				-- Calculate tokens to add based on time elapsed in seconds
				local elapsed = (now - lastRefillTime) / 1e9
				local tokensToAdd = elapsed * refillRate
				currentTokens = math.min(currentTokens + tokensToAdd, maxTokens)
			end

			if currentTokens >= tokensRequested then
				currentTokens = currentTokens - tokensRequested

				redis.call("HMSET", KEYS[1], "lastRefillTime", now, "currentTokens", currentTokens)
				redis.call("PEXPIRE", KEYS[1], 3600000)
				return 1
			else
				-- Not enough tokens, reject the request
				return 0
			end
	`

	now := time.Now().UnixNano()
	result, err := rtb.client.Eval(ctx, script, []string{rtb.bucketKey}, rtb.maxTokens, rtb.refillRate, tokens, now).Int()
	if err != nil {
		return false, fmt.Errorf("error executing Lua script: %w", err)
	}

	return result == 1, nil
}
