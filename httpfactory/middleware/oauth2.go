package middleware

import (
    "context"
    "log"
    "net/http"
    "sync"
    "time"
)

// OAuth2Token holds token info.
type OAuth2Token struct {
    AccessToken  string
    RefreshToken string
    Expiry       time.Time
}

// OAuth2Config holds refresh logic.
type OAuth2Config struct {
    GetTokenFunc func(ctx context.Context, refreshToken string) (OAuth2Token, error)
    Token        OAuth2Token
    mu           sync.Mutex
}

// OAuth2Middleware handles token refresh.
func OAuth2Middleware(cfg *OAuth2Config) func(http.RoundTripper) http.RoundTripper {
    return func(next http.RoundTripper) http.RoundTripper {
        return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
            cfg.mu.Lock()
            if time.Now().After(cfg.Token.Expiry) {
                log.Println("🔄 Refreshing OAuth2 Token...")
                newToken, err := cfg.GetTokenFunc(req.Context(), cfg.Token.RefreshToken)
                if err == nil {
                    cfg.Token = newToken
                }
            }
            req.Header.Set("Authorization", "Bearer "+cfg.Token.AccessToken)
            cfg.mu.Unlock()
            return next.RoundTrip(req)
        })
    }
}

