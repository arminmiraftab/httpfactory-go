package httpfactory

import (
    "io"
    "net/url"
    "strings"
)

func BuildFormBody(data map[string]string) io.Reader {
    form := url.Values{}
    for k, v := range data {
        form.Set(k, v)
    }
    return strings.NewReader(form.Encode())
}