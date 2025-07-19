package main

import (
    "bytes"
    "context"
    "fmt"
    "io"
    "log"
    "mime/multipart"
    "os"
    "time"

    "github.com/arminmiraftab/httpfactory-go/httpfactory"
    "github.com/arminmiraftab/httpfactory-go/httpfactory/middleware"
)

func main() {
    client := httpfactory.NewClient(
        15*time.Second,
        middleware.AuthMiddlewareFactory(middleware.AuthConfig{
            Type:  middleware.AuthBearer,
            Token: "your-bearer-token",
        }),
        middleware.LoggerMiddleware(),
    )

    // 1️⃣ JSON Body
    jsonBody := []byte(`{"title":"Hello","body":"This is a JSON payload."}`)
    resp, err := client.Post(
        context.Background(),
        "https://httpbin.org/post",
        "application/json",
        jsonBody,
        nil,
    )
    if err != nil {
        log.Fatal("JSON POST failed:", err)
    }
    fmt.Println("JSON POST status:", resp.StatusCode)
    fmt.Println(string(resp.Body))

    // 2️⃣ Form URL‑encoded Body
    form := httpfactory.BuildFormURLEncodedBody(map[string]string{
        "username": "alice",
        "password": "s3cr3t",
    })
    // form is an io.Reader, so read its bytes:
    buf := new(bytes.Buffer)
    io.Copy(buf, form)
    resp, err = client.Post(
        context.Background(),
        "https://httpbin.org/post",
        "application/x-www-form-urlencoded",
        buf.Bytes(),
        nil,
    )
    if err != nil {
        log.Fatal("Form POST failed:", err)
    }
    fmt.Println("Form POST status:", resp.StatusCode)
    fmt.Println(string(resp.Body))

    // 3️⃣ Multipart Form‑Data Body (فایل آپلود)
    // آماده‌سازی فیلدها و فایل‌ها
    var b bytes.Buffer
    mw := multipart.NewWriter(&b)

    // فیلد متنی
    mw.WriteField("description", "My upload file")

    // فیلد فایل
    file, err := os.Open("path/to/localfile.jpg")
    if err != nil {
        log.Fatal("Open file failed:", err)
    }
    defer file.Close()

    part, err := mw.CreateFormFile("file", "localfile.jpg")
    if err != nil {
        log.Fatal(err)
    }
    io.Copy(part, file)
    mw.Close()

    resp, err = client.Post(
        context.Background(),
        "https://httpbin.org/post",
        mw.FormDataContentType(),
        b.Bytes(),
        nil,
    )
    if err != nil {
        log.Fatal("Multipart POST failed:", err)
    }
    fmt.Println("Multipart POST status:", resp.StatusCode)
    fmt.Println(string(resp.Body))
}
