package httpfactory

import (
    "io"
    "net/url"
    "strings"
)

// BuildFormURLEncodedBody builds x-www-form-urlencoded body
func BuildFormURLEncodedBody(data map[string]string) io.Reader {
    form := url.Values{}
    for k, v := range data {
        form.Set(k, v)
    }
    return strings.NewReader(form.Encode())
}
