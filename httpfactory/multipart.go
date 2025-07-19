package httpfactory

import (
    "bytes"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
)

func BuildMultipartBody(fields map[string]string, files map[string]string) (io.Reader, string, error) {
    var buf bytes.Buffer
    writer := multipart.NewWriter(&buf)

    for key, val := range fields {
        if err := writer.WriteField(key, val); err != nil {
            return nil, "", err
        }
    }

    for fieldname, filepath := range files {
        file, err := os.Open(filepath)
        if err != nil {
            return nil, "", err
        }
        defer file.Close()

        part, err := writer.CreateFormFile(fieldname, filepath.Base(filepath))
        if err != nil {
            return nil, "", err
        }

        if _, err := io.Copy(part, file); err != nil {
            return nil, "", err
        }
    }

    err := writer.Close()
    if err != nil {
        return nil, "", err
    }

    return &buf, writer.FormDataContentType(), nil
}