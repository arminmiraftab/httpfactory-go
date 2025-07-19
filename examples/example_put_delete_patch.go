package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/arminmiraftab/httpfactory-go/httpfactory"
)

func main() {
    client := httpfactory.NewClient(10*time.Second)

    // PUT example
    putData := []byte("updated content")
    putResp, err := client.Put(context.Background(), "https://httpbin.org/put", "text/plain", putData, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("PUT Status:", putResp.StatusCode)

    // DELETE example
    delResp, err := client.Delete(context.Background(), "https://httpbin.org/delete", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("DELETE Status:", delResp.StatusCode)

    // PATCH example
    patchData := []byte(`{"status":"patched"}`)
    patchResp, err := client.Patch(context.Background(), "https://httpbin.org/patch", "application/json", patchData, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("PATCH Status:", patchResp.StatusCode)
}
