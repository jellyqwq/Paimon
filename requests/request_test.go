package requests

import (
	"testing"
)

func TestBronya(t *testing.T) {
	res, _ := Bronya("POST", "https://httpbin.org/post", nil, `{"test":"ok"}`)
	t.Logf("ResCode: %d", res.StatusCode)
	t.Log(res.Header)
	t.Logf("Content: %d", string(res.Body))
}
