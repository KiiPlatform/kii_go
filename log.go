// +build !debug

package kii

import "net/http"

func logRequest(req *http.Request, reqBody []byte, resp *http.Response, respBody []byte) {
	// nothing to do.
}
