package test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/arminmiraftab/httpfactory-go/httpfactory"
    "github.com/arminmiraftab/httpfactory-go/httpfactory/middleware"
)

func TestRetryMiddleware(t *testing.T) {
    attempts := 0
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        attempts++
        if attempts < 2 {
            http.Error(w, "error", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusOK)
    }))
    defer ts.Close()

    retryCfg := middleware.RetryConfig{MaxRetries: 1, BaseDelay: 10 * time.Millisecond, MaxDelay: 20 * time.Millisecond, EnableJitter: false}
    client := httpfactory.NewClient(1*time.Second, middleware.RetryMiddleware(retryCfg))
    resp, err := client.Get(context.Background(), ts.URL, nil)
    if err != nil {
        t.Fatalf("expected success, got %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected 200, got %d", resp.StatusCode)
    }
}

