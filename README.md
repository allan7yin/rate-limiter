![SafetyGopher](./docs/file.png)
# GopherGate - Simple Rate Limiter in Go

---

A simple and reusable request rate-limiting microservice that protects external applications' endpoints.
## Idea
```go
const (
	MAX_BUCKET_SIZE float64 = 10
	REFILL_RATE     int     = 1
)

type TokenBucket struct {
	currentBucketSize   float64
	lastRefillTimestamp int64
}

func (tb *TokenBucket) allowRequest(tokens float64) bool {
	tb.refill()

	if tb.currentBucketSize >= tokens {
		tb.currentBucketSize -= tokens
		return true
	}

	return false
}

func getCurrentTimeInNanoseconds() int64 {
	return time.Now().UnixNano()
}

func (tb *TokenBucket) refill() {
	nowTime := getCurrentTimeInNanoseconds()
	elapsedTime := nowTime - tb.lastRefillTimestamp
	tokensToAdd := float64(elapsedTime) * float64(REFILL_RATE) / 1e9 * 2
	tb.currentBucketSize = math.Min(tb.currentBucketSize+tokensToAdd, MAX_BUCKET_SIZE)
	tb.lastRefillTimestamp = nowTime
}
```

Where we make use `Redis` as a token bucket due to its performance in distributed settings.
