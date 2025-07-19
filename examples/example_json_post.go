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
        middleware.LoggerMiddleware(),
    )

    jsonBody := []byte(`{"name":"Alice","age":30}`)
    resp, err := client.Post(
        context.Background(),
        "https://httpbin.org/post",
        "application/json",
        jsonBody,
        nil,
    )
    if err != nil {
        log.Fatal(err)	
    }
    fmt.Println("POST Status:", resp.StatusCode)
    fmt.Println("Response:", string(resp.Body))
}
