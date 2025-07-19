package middleware

import (
    "errors"
    "net/http"
    "sync"
    "time"
)

type CircuitBreakerConfig struct {
    FailureThreshold int
    RetryTimeout     time.Duration
}

type circuitState int

const (
    stateClosed circuitState = iota
    stateOpen
    stateHalfOpen
)

type CircuitBreaker struct {
    cfg         CircuitBreakerConfig
    state       circuitState
    failures    int
    lastFailure time.Time
    mu          sync.Mutex
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
    return &CircuitBreaker{cfg: cfg}
}

// Middleware returns the circuit breaker middleware.
func (cb *CircuitBreaker) Middleware() func(http.RoundTripper) http.RoundTripper {
    return func(next http.RoundTripper) http.RoundTripper {
        return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
            cb.mu.Lock()
            if cb.state == stateOpen && time.Since(cb.lastFailure) < cb.cfg.RetryTimeout {
                cb.mu.Unlock()
                return nil, errors.New("circuit breaker open")
            }
            if cb.state == stateOpen {
                cb.state = stateHalfOpen
            }
            cb.mu.Unlock()

            resp, err := next.RoundTrip(req)

            cb.mu.Lock()
            defer cb.mu.Unlock()
            if err != nil || resp.StatusCode >= 500 {
                cb.failures++
                cb.lastFailure = time.Now()
                if cb.failures >= cb.cfg.FailureThreshold {
                    cb.state = stateOpen
                }
            } else {
                cb.failures = 0
                if cb.state == stateHalfOpen {
                    cb.state = stateClosed
                }
            }
            return resp, err
        })
    }
}

