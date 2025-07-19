package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/arminmiraftab/httpfactory-go/httpfactory"
    "github.com/arminmiraftab/httpfactory-go/httpfactory/middleware"
)

func main() {
    client := httpfactory.NewClient(
        10*time.Second,
        middleware.AuthMiddlewareFactory(middleware.AuthConfig{
            Type:  middleware.AuthBearer,
            Token: "your-bearer-token",
        }),
        middleware.LoggerMiddleware(),
    )

    resp, err := client.Get(context.Background(), "https://httpbin.org/get", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Status:", resp.StatusCode)
    fmt.Println("Body:", string(resp.Body))
}
