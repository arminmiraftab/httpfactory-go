package httpfactory

import "net/http"

// CustomHeadersMiddleware injects custom headers
func CustomHeadersMiddleware(headers map[string]string) Middleware {
    return func(next http.RoundTripper) http.RoundTripper {
        return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
            for k, v := range headers {
                req.Header.Set(k, v)
            }
            return next.RoundTrip(req)
        })
    }
}
