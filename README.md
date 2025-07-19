# httpfactory-go

Optimized HTTP Client for Go with pluggable middleware and tests.

## Features

- HTTP methods: GET, POST, PUT, PATCH, DELETE
- Middleware: Logger, Retry (backoff+jitter), CircuitBreaker, Auth (Bearer/Basic/APIKey/OAuth2), CustomHeaders
- Build x-www-form-urlencoded bodies
- Fully modular and testable

## Usage

```go
import (
    "context"
    "time"

    "github.com/arminmiraftab/httpfactory-go/httpfactory"
    "github.com/arminmiraftab/httpfactory-go/httpfactory/middleware"
)

func main() {
    client := httpfactory.NewClient(10*time.Second,
        middleware.AuthMiddlewareFactory(middleware.AuthConfig{
            Type:  middleware.AuthBearer,
            Token: "token123",
        }),
        middleware.LoggerMiddleware(),
        middleware.RetryMiddleware(middleware.RetryConfig{
            MaxRetries:   3,
            BaseDelay:    100 * time.Millisecond,
            MaxDelay:     1 * time.Second,
            EnableJitter: true,
        }),
    )

    resp, err := client.Get(context.Background(), "https://httpbin.org/get", nil)
    if err != nil {
        panic(err)
    }
    fmt.Println("Status:", resp.StatusCode)
    fmt.Println("Body:", string(resp.Body))
}
```
