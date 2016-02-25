// +build debug

package kii

import (
	"bytes"
	"fmt"
	"net/http"
)

func headerToString(h http.Header) string {
	var (
		b        = new(bytes.Buffer)
		notFirst bool
	)
	b.WriteString("{")
	for k, v := range h {
		if notFirst {
			b.WriteString(",")
		} else {
			notFirst = true
		}
		fmt.Fprintf(b, "%q:", k)
		switch len(v) {
		case 0:
			b.WriteString("(nil)")
		case 1:
			fmt.Fprintf(b, "%q", v[0])
		default:
			fmt.Fprintf(b, "%#v", v)
		}
	}
	b.WriteString("}")
	return b.String()
}

// logRequest logs request and response.
func logRequest(req *http.Request, resp *http.Response, respBody []byte) {
	var s1, s2 string
	// TODO: setup s1.
	if len(respBody) > 0 {
		s2 = string(respBody)
	}
	Logger.Printf(`access to Kii:
  url=%s
  method=%s
  header=%s
  body=%q
  status_code=%d
  response_header=%s
  response_body=%q`,
		req.URL,
		req.Method,
		headerToString(req.Header),
		s1,
		resp.StatusCode,
		headerToString(resp.Header),
		s2)
}
