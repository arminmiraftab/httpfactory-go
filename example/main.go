package main

import (
    "fmt"
    "log"
    "time"

    "github.com/arminmiraftab/httpfactory-go/httpfactory"
)

type simpleLogger struct{}

func (l *simpleLogger) Debug(msg string) { log.Println("[DEBUG]", msg) }
func (l *simpleLogger) Info(msg string)  { log.Println("[INFO]", msg) }
func (l *simpleLogger) Error(msg string) { log.Println("[ERROR]", msg) }

func main() {
    client := httpfactory.NewClient(httpfactory.Config{
        Timeout:       10 * time.Second,
        RetryCount:    3,
        RetryInterval: 2 * time.Second,
        DefaultHeaders: map[string]string{
            "Authorization": "Bearer token123",
        },
        Logger: &simpleLogger{},
    })

    // JSON POST Example
    jsonBody, err := httpfactory.BuildJSONBody(map[string]interface{}{
        "name": "John Doe",
        "age":  30,
    })
    if err != nil {
        panic(err)
    }
    resp, err := client.Post("https://httpbin.org/post", "application/json", jsonBody, nil)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    fmt.Println("POST JSON Status:", resp.Status)
}