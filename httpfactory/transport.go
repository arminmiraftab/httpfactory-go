package httpfactory

import "net/http"

// Chain applies middleware in order
func Chain(rt http.RoundTripper, mws ...Middleware) http.RoundTripper {
    if rt == nil {
        rt = http.DefaultTransport
    }
    for _, mw := range mws {
        rt = mw(rt)
    }
    return rt
}

// roundTripperFunc adapter
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
    return f(req)
}
