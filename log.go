// +build !debug

package kii

import "net/http"

func logRequest(req *http.Request, resp *http.Response, body []byte) {
	// nothing to do.
}
