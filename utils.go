package requests

import (
	"net/http"
	"net/url"
)

const (
	Version = "1.0"
)

var (
	sessionWithoutCookies *session
)

func init() {
	client := &http.Client{}
	s := &session{Client: client}
	s.cookies = make(map[string]string)
	sessionWithoutCookies = s
}

func parseParams(params map[string][]string) url.Values {
	v := url.Values{}
	for key, values := range params {
		for _, value := range values {
			v.Add(key, value)
		}
	}
	return v
}

func parseHeaders(headers map[string][]string) http.Header {
	h := http.Header{}
	for key, values := range headers {
		for _, value := range values {
			h.Add(key, value)
		}
	}
	_, hasAccept := h["Accept"]
	if !hasAccept {
		h.Add("Accept", "*/*")
	}
	_, hasAgent := h["User-Agent"]
	if !hasAgent {
		h.Add("User-Agent", "go-requests/"+Version)
	}
	return h
}

func Get(urlPath string) (Request, error) {
	return newRequest("GET", urlPath, sessionWithoutCookies)
}
