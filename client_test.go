package kii

import "testing"

func TestSetDefaultUserAgent(t *testing.T) {
	checkRequest := func(expectedUA string) {
		r, err := newRequest("GET", "http://example.com", nil)
		if err != nil {
			t.Errorf("newRequest() faild: %s", err)
			return
		}
		ua := r.UserAgent()
		if ua != expectedUA {
			t.Errorf("UA not matched: %q (expected: %q)", ua, expectedUA)
		}
	}

	if defaultUserAgent != "" {
		t.Fatal("initial value of defaultUserAgent should be empty")
	}
	checkRequest("")

	SetDefaultUserAgent("foo")
	if defaultUserAgent != "foo" {
		t.Fatalf("defaultUserAgent should be %q but %q", "foo", defaultUserAgent)
	}
	checkRequest("foo")

	SetDefaultUserAgent("bar")
	if defaultUserAgent != "bar" {
		t.Fatalf("defaultUserAgent should be %q but %q", "bar", defaultUserAgent)
	}
	checkRequest("bar")

	SetDefaultUserAgent("")
	if defaultUserAgent != "" {
		t.Fatalf("defaultUserAgent should be empty but %q", defaultUserAgent)
	}
	checkRequest("")
}
