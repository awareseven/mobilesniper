package utils

import (
	"bytes"
	"io"
	"net/http"
)

// Do an HTTP SOAP request on given URL and data.
// The requests header will be set automatically to mock a sdcX client.
func DoSdcXClientSOAPPost(url string, soapData []byte, verbose ...bool) (*http.Response, []byte, error) {
	var v bool
	if verbose == nil {
		v = false
	} else {
		v = verbose[0]
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(soapData))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("User-Agent", "sdcX/0.1")
	req.Header.Set("Content-Type", "application/soap+xml")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	LogVerbosef(v, "Received %d bytes response:\n%s", resp.ContentLength, respBytes)

	return resp, respBytes, err
}
