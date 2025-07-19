# httpfactory-go

Enterprise-grade HTTP Client for Go with pluggable Middleware.

## Features

- HTTP methods: GET, POST
- Middleware: Logger, Retry, CircuitBreaker, Auth, APIKey, CustomHeaders
- Auth types: Bearer, Basic, API Key, OAuth2 (auto-refresh), Custom
- Retry with exponential backoff + jitter
- Circuit Breaker pattern
- Build x-www-form-urlencoded bodies
- Fully modular and testable

## Usage

```go
client := httpfactory.New(
    10*time.Second,
    httpfactory.AuthMiddlewareFactory(httpfactory.AuthConfig{
        Type:          httpfactory.AuthBearer,
        Token:         "token123",
    }),
    httpfactory.LoggerMiddleware(),
)
```
