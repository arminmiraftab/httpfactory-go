package httpfactory

import "net/http"

// APIKeyMiddleware sets API key header
func APIKeyMiddleware(headerName, apiKey string) Middleware {
    return func(next http.RoundTripper) http.RoundTripper {
        return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
            req.Header.Set(headerName, apiKey)
            return next.RoundTrip(req)
        })
    }
}
