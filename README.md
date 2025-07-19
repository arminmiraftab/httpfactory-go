# httpfactory-go

یک کلاینت HTTP حرفه‌ای و ماژولار برای زبان Go  
پشتیبانی از انواع درخواست‌ها و انواع داده‌ها با قابلیت Retry و Timeout

## امکانات
- متدهای HTTP کامل: GET, POST, PUT, PATCH, DELETE
- ارسال JSON, فرم دیتا و multipart فایل
- مدیریت Retry هوشمند و Timeout
- Default Headers و Logging
- طراحی ماژولار و قابل توسعه

## شروع سریع

```go
client := httpfactory.NewClient(httpfactory.Config{
    Timeout:       10 * time.Second,
    RetryCount:    3,
    RetryInterval: 2 * time.Second,
    DefaultHeaders: map[string]string{
        "Authorization": "Bearer token",
    },
    Logger: &simpleLogger{},
})

jsonBody, _ := httpfactory.BuildJSONBody(map[string]interface{}{
    "name": "John",
})

resp, err := client.Post("https://api.example.com/users", "application/json", jsonBody, nil)
```