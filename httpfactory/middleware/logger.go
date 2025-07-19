package middleware

import (
    "log"
    "net/http"
    "time"
)

// LoggerMiddleware logs requests and responses.
func LoggerMiddleware() func(http.RoundTripper) http.RoundTripper {
    return func(next http.RoundTripper) http.RoundTripper {
        return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
            start := time.Now()
            log.Printf("➡️ %s %s", req.Method, req.URL)
            resp, err := next.RoundTrip(req)
            if err != nil {
                log.Printf("❌ %v", err)
                return nil, err
            }
            log.Printf("⬅️ %d %s (%v)", resp.StatusCode, req.URL, time.Since(start))
            return resp, nil
        })
    }
}

