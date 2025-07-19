package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "net/url"

    "github.com/arminmiraftab/httpfactory-go/httpfactory"
)

func main() {
    oauthCfg := &httpfactory.OAuth2Config{
        Token: httpfactory.OAuth2Token{
            AccessToken:  "initial",
            RefreshToken: "refresh",
            Expiry:       time.Now().Add(-1 * time.Minute),
        },
        GetTokenFunc: func(ctx context.Context, refreshToken string) (httpfactory.OAuth2Token, error) {
            log.Println("Fetching new token...")
            return httpfactory.OAuth2Token{
                AccessToken:  "new-token",
                RefreshToken: refreshToken,
                Expiry:       time.Now().Add(1 * time.Hour),
            }, nil
        },
    }

    client := httpfactory.New(
        10*time.Second,
        httpfactory.AuthMiddlewareFactory(httpfactory.AuthConfig{
            Type:          httpfactory.AuthOAuth2,
            OAuth2Cfg:     oauthCfg,
            CustomHeaders: map[string]string{"X-Custom": "value"},
            CustomParams:  url.Values{"debug": {"true"}},
        }),
        httpfactory.LoggerMiddleware(),
        httpfactory.RetryMiddleware(httpfactory.RetryConfig{
            MaxRetries:   3,
            BaseDelay:    200 * time.Millisecond,
            MaxDelay:     2 * time.Second,
            EnableJitter: true,
        }),
        httpfactory.APIKeyMiddleware("X-API-KEY", "apikey123"),
        httpfactory.CustomHeadersMiddleware(map[string]string{"X-Another": "val"}),
    )

    resp, err := client.Get("https://httpbin.org/get", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Status:", resp.StatusCode)
    fmt.Println("Headers:", resp.Headers)
    fmt.Println("Body:", string(resp.Body))
}
