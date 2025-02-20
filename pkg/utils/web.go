package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func IsHTTPorHTTPS(ip string, port int) bool {
	urlHTTP := fmt.Sprintf("http://%s:%d", ip, port)
	urlHTTPS := fmt.Sprintf("https://%s:%d", ip, port)

	// Check HTTP
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(urlHTTP)
	if err == nil {
		resp.Body.Close()
		return true
	}

	// Check HTTPS
	client = &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err = client.Get(urlHTTPS)
	if err == nil {
		resp.Body.Close()
		return true
	}

	return false
}

func IsHTTPPayload(payload []byte) bool {
	// List of common HTTP methods
	httpMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "CONNECT"}

	for _, method := range httpMethods {
		if bytes.HasPrefix(payload, []byte(method+" ")) {
			return true
		}
	}

	// Check for HTTP/1.1 or HTTP/2.0 in response lines
	if bytes.Contains(payload, []byte("HTTP/1.1")) || bytes.Contains(payload, []byte("HTTP/2.0")) {
		return true
	}

	return false
}
