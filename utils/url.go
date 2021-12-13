package utils

import (
    "fmt"
    "net/url"
)

type URL struct {
    url.URL
}

func ParseURL(s string) (*URL, error) {
    u, err := url.Parse(s)
    if err != nil {
        return nil, err
    }
    if u.Scheme != "http" && u.Scheme != "https" {
        return nil, fmt.Errorf("unsupported scheme %q for URL", u.Scheme)
    }
    if u.Host == "" {
        return nil, fmt.Errorf("missing host for URL")
    }
    return &URL{*u}, nil
}
