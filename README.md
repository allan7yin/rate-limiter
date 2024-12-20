![SafetyGopher](./docs/file.png)
# GopherGate - Simple Rate Limiter in Go

A simple and reusable request rate-limiting microservice that protects external applications' endpoints.

---
## Performance
This was designed to handle just under 1,000,000 requests per minute. To tune for this, the service is designed to handle ~16,667 requests
per second. Maximum bucket size is how many requests that can be handled under bursty conditions, which this service sets as 
5 seconds. So, refill rate set to 16,667 tokens per second, with a maximum bucket size of 83, 335 tokens. 

To simulate these conditions in Jmeter, the following configuration was used:
- 500 threads
- 10 second ramp up time
- 2000 loop count

The results of this are:

<img src="./docs/Base.png" alt="BaseMetric" width="300">

## Idea
This limiter uses `Redis` as a token bucket due to its performance in distributed settings. Below is a simple illustration of what the system is like:
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
