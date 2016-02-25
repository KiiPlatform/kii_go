package kii

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type contentTyper interface {
	contentType() string
}

// newRequest creates http.Request with JSON body and header.
func newRequest(method, url string, body interface{}) (*http.Request, error) {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}
	// set Content-Type if available automatically.
	if body != nil {
		if t, ok := body.(contentTyper); ok {
			req.Header.Set("Content-Type", t.contentType())
		} else {
			req.Header.Set("Content-Type", "application/json")
		}
	}
	return req, nil
}

func executeRequest(req *http.Request) ([]byte, error) {
	return executeRequest2(req, 200, 400)
}

func executeRequest2(req *http.Request, scMin, scMax int) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logRequest(req, resp, b)

	if resp.StatusCode < scMin || resp.StatusCode >= scMax {
		return nil, errors.New(string(b))
	}
	return b, nil
}
