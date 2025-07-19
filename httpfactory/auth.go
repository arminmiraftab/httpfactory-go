package httpfactory

import (
    "encoding/base64"
    "io"
    "net/http"
    "net/url"
    "strings"
    "time"
)

// AuthType enum-like for auth types
type AuthType int

const (
    AuthNone AuthType = iota
    AuthBearer
    AuthBasic
    AuthAPIKey
    AuthOAuth2
    AuthCustom
)

// AuthConfig holds all auth options
type AuthConfig struct {
    Type          AuthType
    Token         string
    Username      string
    Password      string
    APIKeyHeader  string
    OAuth2Cfg     *OAuth2Config
    CustomHeaders map[string]string
    CustomParams  url.Values
    CustomBody    string
}

// AuthMiddlewareFactory returns auth middleware based on type
func AuthMiddlewareFactory(cfg AuthConfig) Middleware {
    return func(next http.RoundTripper) http.RoundTripper {
        return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
            if len(cfg.CustomParams) > 0 {
                q := req.URL.Query()
                for k, vals := range cfg.CustomParams {
                    for _, v := range vals {
                        q.Add(k, v)
                    }
                }
                req.URL.RawQuery = q.Encode()
            }
            for k, v := range cfg.CustomHeaders {
                req.Header.Set(k, v)
            }
            switch cfg.Type {
            case AuthBearer:
                if cfg.Token != "" {
                    req.Header.Set("Authorization", "Bearer "+cfg.Token)
                }
            case AuthBasic:
                if cfg.Username != "" && cfg.Password != "" {
                    auth := base64.StdEncoding.EncodeToString([]byte(cfg.Username + ":" + cfg.Password))
                    req.Header.Set("Authorization", "Basic "+auth)
                }
            case AuthAPIKey:
                if cfg.APIKeyHeader != "" && cfg.Token != "" {
                    req.Header.Set(cfg.APIKeyHeader, cfg.Token)
                }
            case AuthOAuth2:
                if cfg.OAuth2Cfg != nil {
                    cfg.OAuth2Cfg.mu.Lock()
                    if time.Now().After(cfg.OAuth2Cfg.Token.Expiry) {
                        newToken, err := cfg.OAuth2Cfg.GetTokenFunc(req.Context(), cfg.OAuth2Cfg.Token.RefreshToken)
                        if err == nil {
                            cfg.OAuth2Cfg.Token = newToken
                        }
                    }
                    req.Header.Set("Authorization", "Bearer "+cfg.OAuth2Cfg.Token.AccessToken)
                    cfg.OAuth2Cfg.mu.Unlock()
                }
            case AuthCustom:
                if cfg.CustomBody != "" {
                    req.Body = io.NopCloser(strings.NewReader(cfg.CustomBody))
                }
            }
            return next.RoundTrip(req)
        })
    }
}
