package httpfactory

import (
    "bytes"
    "context"
    "io"
    "net/http"
    "time"
)

// Client is the HTTP client with middleware support.
type Client struct {
    http.Client
}

// Middleware defines a RoundTripper decorator.
type Middleware func(next http.RoundTripper) http.RoundTripper

// NewClient creates a new Client with the given timeout and middleware chain.
func NewClient(timeout time.Duration, mws ...Middleware) *Client {
    transport := Chain(http.DefaultTransport, mws...)
    return &Client{Client: http.Client{
        Timeout:   timeout,
        Transport: transport,
    }}
}

// Response wraps the http.Response data.
type Response struct {
    StatusCode int
    Headers    http.Header
    Body       []byte
}

// do executes an HTTP request with context, body reader and headers.
func (c *Client) do(ctx context.Context, method, url string, body io.Reader, headers map[string]string) (*Response, error) {
    req, err := http.NewRequestWithContext(ctx, method, url, body)
    if err != nil {
        return nil, err
    }
    for k, v := range headers {
        req.Header.Set(k, v)
    }
    resp, err := c.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    return &Response{resp.StatusCode, resp.Header, data}, nil
}

// Get sends a GET request.
func (c *Client) Get(ctx context.Context, url string, headers map[string]string) (*Response, error) {
    return c.do(ctx, http.MethodGet, url, nil, headers)
}

// Post sends a POST request with content-type header.
func (c *Client) Post(ctx context.Context, url, contentType string, body []byte, headers map[string]string) (*Response, error) {
    if headers == nil {
        headers = map[string]string{}
    }
    headers["Content-Type"] = contentType
    return c.do(ctx, http.MethodPost, url, bytes.NewReader(body), headers)
}

func (c *Client) Put(ctx context.Context, url, contentType string, body []byte, headers map[string]string) (*Response, error) {
    if headers == nil {
        headers = map[string]string{}
    }
    headers["Content-Type"] = contentType
    return c.do(ctx, http.MethodPut, url, bytes.NewReader(body), headers)
}
// Delete sends a DELETE request.
func (c *Client) Delete(ctx context.Context, url string, headers map[string]string) (*Response, error) {
    return c.do(ctx, http.MethodDelete, url, nil, headers)
}

// Patch sends a PATCH request with content-type header.
func (c *Client) Patch(ctx context.Context, url, contentType string, body []byte, headers map[string]string) (*Response, error) {
    if headers == nil {
        headers = map[string]string{}
    }
    headers["Content-Type"] = contentType
    return c.do(ctx, http.MethodPatch, url, bytes.NewReader(body), headers)
}

