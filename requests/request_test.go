package requests

import (
	"fmt"
	"testing"
)

func TestBronya(t *testing.T) {
	text := `moon hola`
	data := fmt.Sprintf(`{
		"context": {
			"client": {
				"clientName": "WEB",
				"clientVersion": "2.20220617.00.00"
			}
		},
		"query": "%v"
	}`, text)
	t.Logf("data: %v", data)
	res, _ := Bronya("POST", "https://www.youtube.com/youtubei/v1/search", nil, data)
	t.Logf("ResCode: %d", res.StatusCode)
	t.Log(res.Header)
	t.Logf("Content: %s", string(res.Body))
}
