package test

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/arminmiraftab/httpfactory-go/httpfactory"
)

func TestClient_Get_Succeeds(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Test", "value")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("hello"))
    }))
    defer ts.Close()

    client := httpfactory.NewClient(2*time.Second)
    resp, err := client.Get(context.Background(), ts.URL, nil)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("unexpected status: %d", resp.StatusCode)
    }
    if string(resp.Body) != "hello" {
        t.Errorf("unexpected body: %s", resp.Body)
    }
    if resp.Headers.Get("X-Test") != "value" {
        t.Errorf("missing header X-Test")
    }
}

