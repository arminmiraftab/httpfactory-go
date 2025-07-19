package httpfactory

import (
    "bytes"
    "encoding/json"
    "io"
)

func BuildJSONBody(data interface{}) (io.Reader, error) {
    b, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }
    return bytes.NewReader(b), nil
}