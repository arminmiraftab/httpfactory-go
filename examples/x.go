package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "bytes"
    "net/url"

    "github.com/arminmiraftab/httpfactory-go/httpfactory"
    "github.com/arminmiraftab/httpfactory-go/httpfactory/middleware"
)

func main() {
    // ➊ پیکربندی Middleware: 
    //    - اضافه کردن کوئری پارامتر با CustomParams
    //    - لاگ
    client := httpfactory.NewClient(
        15*time.Second,
        middleware.AuthMiddlewareFactory(middleware.AuthConfig{
            Type:         middleware.AuthCustom,
            CustomParams: url.Values{"search": {"golang"}, "page": {"2"}},
        }),
        middleware.LoggerMiddleware(),
    )

    // ➋ JSON Body
    body := []byte(`{"filter":"recent","limit":10}`)

    // ➌ ساخت URL همراه کوئری
    //    در واقع Middleware این را اضافه می‌کند، 
    //    اما برای نمایش می‌توانید چاپ کنید:
    rawURL := "https://httpbin.org/anything"
    fmt.Println("Request URL (with query):", rawURL+"?search=golang&page=2")

    // ➍ ارسال POST
    resp, err := client.Post(
        context.Background(),
        rawURL,
        "application/json",
        body,
        map[string]string{"X-Custom":"hello"},
    )
    if err != nil {
        log.Fatalf("Request failed: %v", err)
    }

    // ➎ نمایش پاسخ
    fmt.Println("Status:", resp.StatusCode)
    fmt.Println("Response Headers:", resp.Headers)
    fmt.Println("Response Body:", string(resp.Body))
}
