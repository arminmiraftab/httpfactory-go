package httpfactory

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "time"
)

type HttpClient interface {
    Do(req *http.Request) (*http.Response, error)
    Get(url string, headers map[string]string) (*http.Response, error)
    Post(url string, contentType string, body io.Reader, headers map[string]string) (*http.Response, error)
    Put(url string, contentType string, body io.Reader, headers map[string]string) (*http.Response, error)
    Patch(url string, contentType string, body io.Reader, headers map[string]string) (*http.Response, error)
    Delete(url string, headers map[string]string) (*http.Response, error)
}

type Config struct {
    Timeout        time.Duration
    RetryCount     int
    RetryInterval  time.Duration
    DefaultHeaders map[string]string
    Logger         Logger
}

type Logger interface {
    Debug(msg string)
    Info(msg string)
    Error(msg string)
}

type httpClient struct {
    client         *http.Client
    retryCount     int
    retryInterval  time.Duration
    defaultHeaders map[string]string
    logger         Logger
}

func NewClient(cfg Config) HttpClient {
    return &httpClient{
        client: &http.Client{
            Timeout: cfg.Timeout,
        },
        retryCount:     cfg.RetryCount,
        retryInterval:  cfg.RetryInterval,
        defaultHeaders: cfg.DefaultHeaders,
        logger:         cfg.Logger,
    }
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
    for k, v := range c.defaultHeaders {
        if req.Header.Get(k) == "" {
            req.Header.Set(k, v)
        }
    }

    var resp *http.Response
    var err error
    for attempt := 0; attempt <= c.retryCount; attempt++ {
        if c.logger != nil {
            c.logger.Debug(fmt.Sprintf("HTTP Request: %s %s Attempt: %d", req.Method, req.URL.String(), attempt))
        }

        resp, err = c.client.Do(req)
        if err == nil && resp.StatusCode < 500 {
            break
        }

        if attempt < c.retryCount {
            if c.logger != nil {
                c.logger.Info("Retrying in " + c.retryInterval.String())
            }
            time.Sleep(c.retryInterval)
        }
    }
    return resp, err
}

func (c *httpClient) Get(url string, headers map[string]string) (*http.Response, error) {
    return c.doRequest("GET", url, "", nil, headers)
}

func (c *httpClient) Post(url string, contentType string, body io.Reader, headers map[string]string) (*http.Response, error) {
    return c.doRequest("POST", url, contentType, body, headers)
}

func (c *httpClient) Put(url string, contentType string, body io.Reader, headers map[string]string) (*http.Response, error) {
    return c.doRequest("PUT", url, contentType, body, headers)
}

func (c *httpClient) Patch(url string, contentType string, body io.Reader, headers map[string]string) (*http.Response, error) {
    return c.doRequest("PATCH", url, contentType, body, headers)
}

func (c *httpClient) Delete(url string, headers map[string]string) (*http.Response, error) {
    return c.doRequest("DELETE", url, "", nil, headers)
}

func (c *httpClient) doRequest(method, url, contentType string, body io.Reader, headers map[string]string) (*http.Response, error) {
    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }
    if contentType != "" {
        req.Header.Set("Content-Type", contentType)
    }
    for k, v := range headers {
        req.Header.Set(k, v)
    }
    return c.Do(req)
}