package requests

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	StatusCode int
	Header     map[string][]string
	Body       []byte
}

// Post and Get methods
func Bronya(method string, url string, headers map[string]string, data string) (*Response, error) {
	var request *http.Request
	var err error
	if method == "POST" {
		request, err = http.NewRequest("POST", url, strings.NewReader(data))
	} else if method == "GET" {
		// GET请求一般不携带内容
		request, err = http.NewRequest("GET", url, nil)
	} else {
		return nil, fmt.Errorf("method must be POST or GET")
	}

	if err != nil {
		return nil, err
	}

	// set request header
	for key, val := range headers {
		request.Header.Set(key, val)
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	response := &Response{}
	response.StatusCode = res.StatusCode
	response.Header = res.Header

	defer res.Body.Close()
	response.Body, _ = io.ReadAll(res.Body)

	return response, nil
}
