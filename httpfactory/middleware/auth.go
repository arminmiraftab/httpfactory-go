package middleware

import (
    "encoding/base64"
    "net/http"
    "net/url"
)

// AuthType enum-like for auth types.
type AuthType int

const (
    AuthNone AuthType = iota
    AuthBearer
    AuthBasic
    AuthAPIKey
    AuthCustom
)

// AuthConfig holds all auth options.
type AuthConfig struct {
    Type          AuthType
    Token         string
    Username      string
    Password      string
    APIKeyHeader  string
    CustomHeaders map[string]string
    CustomParams  url.Values
    CustomBody    string
}

// AuthMiddlewareFactory returns auth middleware based on type.
func AuthMiddlewareFactory(cfg AuthConfig) func(http.RoundTripper) http.RoundTripper {
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
                    auth := base64.StdEncoding.EncodeToString([]byte(cfg.Username+":"+cfg.Password))
                    req.Header.Set("Authorization", "Basic "+auth)
                }
            case AuthAPIKey:
                if cfg.APIKeyHeader != "" && cfg.Token != "" {
                    req.Header.Set(cfg.APIKeyHeader, cfg.Token)
                }
            }
            return next.RoundTrip(req)
        })
    }
}

