package httpfactory

import (
    "bytes"
    "io"
    "net/http"
    "time"
)

// HTTPClient encapsulates http.Client with pluggable middleware
type HTTPClient struct {
    client *http.Client
}

// Middleware defines a RoundTripper decorator
type Middleware func(next http.RoundTripper) http.RoundTripper

// New creates a HTTPClient with given timeout and middleware chain
func New(timeout time.Duration, mws ...Middleware) *HTTPClient {
    transport := http.DefaultTransport
    for _, mw := range mws {
        transport = mw(transport)
    }
    return &HTTPClient{
        client: &http.Client{
            Timeout:   timeout,
            Transport: transport,
        },
    }
}

// ResponseWithHeaders holds full response data
type ResponseWithHeaders struct {
    StatusCode int
    Headers    http.Header
    Body       []byte
}

// sendRequest does the actual HTTP request
func (hc *HTTPClient) sendRequest(method, url string, body io.Reader, headers map[string]string) (*ResponseWithHeaders, error) {
    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }
    for k, v := range headers {
        req.Header.Set(k, v)
    }
    resp, err := hc.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    return &ResponseWithHeaders{StatusCode: resp.StatusCode, Headers: resp.Header, Body: respBody}, nil
}

// Get convenience method
func (hc *HTTPClient) Get(url string, headers map[string]string) (*ResponseWithHeaders, error) {
    return hc.sendRequest(http.MethodGet, url, nil, headers)
}

// Post convenience method
func (hc *HTTPClient) Post(url, contentType string, body []byte, headers map[string]string) (*ResponseWithHeaders, error) {
    h := headers
    if h == nil {
        h = map[string]string{}
    }
    h["Content-Type"] = contentType
    return hc.sendRequest(http.MethodPost, url, bytes.NewReader(body), h)
}
