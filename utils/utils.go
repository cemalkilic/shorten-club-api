package utils

import (
    "fmt"
)

func GetFullHTTPUrl(host string, uri string, isSecure bool) string {
    scheme := "http"
    if isSecure {
        scheme = "https"
    }
    return fmt.Sprintf("%s://%s/%s", scheme, host, uri)
}
