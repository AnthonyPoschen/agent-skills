package example

import (
	"net/http"
	"time"
)

func clients(req *http.Request) {
	// ruleid: no-default-http-client
	_, _ = http.DefaultClient.Do(req)
	// ruleid: no-default-http-client
	_, _ = http.Get("https://example.com")
	// ruleid: no-default-http-client
	c := &http.Client{}
	_ = c
	// ok: no-default-http-client
	ok := &http.Client{Timeout: 10 * time.Second}
	_ = ok
}
