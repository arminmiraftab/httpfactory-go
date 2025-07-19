package httpfactory

import (
    "log"
    "math"
    "math/rand"
    "net/http"
    "time"
)

// RetryConfig controls retry behavior
type RetryConfig struct {
    MaxRetries   int
    BaseDelay    time.Duration
    MaxDelay     time.Duration
    EnableJitter bool
}

// RetryMiddleware implements retry with exponential backoff and jitter
func RetryMiddleware(cfg RetryConfig) Middleware {
    return func(next http.RoundTripper) http.RoundTripper {
        return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
            var resp *http.Response
            var err error
            for i := 0; i <= cfg.MaxRetries; i++ {
                resp, err = next.RoundTrip(req)
                if err == nil && resp.StatusCode < 500 {
                    return resp, err
                }
                if i == cfg.MaxRetries {
                    break
                }
                delay := calculateBackoff(cfg.BaseDelay, i, cfg.MaxDelay, cfg.EnableJitter)
                log.Printf("🔁 retry %d %s %s after %v", i+1, req.Method, req.URL, delay)
                time.Sleep(delay)
            }
            return resp, err
        })
    }
}

// calculateBackoff calculates backoff delay with optional jitter
func calculateBackoff(base time.Duration, attempt int, max time.Duration, jitter bool) time.Duration {
    d := base * time.Duration(math.Pow(2, float64(attempt)))
    if d > max {
        d = max
    }
    if jitter {
        j := time.Duration(rand.Int63n(int64(d / 2)))
        return d + j
    }
    return d
}
