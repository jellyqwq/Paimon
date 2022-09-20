package requests

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	StatusCode int
	Header map[string][]string
	Body []byte
}

// Post and Get methods
func Bronya(method string, url string, headers map[string]string, data string) (*Response, error) {
	var request *http.Request
	var err error
	if method == "POST" {
		request, err = http.NewRequest("POST", url, strings.NewReader(data))
	} else if method == "GET" {
		request, err = http.NewRequest("POST", url, strings.NewReader(data))
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

	// build http client
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	
	var response *Response = new(Response)
	(*response).StatusCode = res.StatusCode
	(*response).Header = res.Header

	defer res.Body.Close()
	(*response).Body, _ = io.ReadAll(res.Body)
	
	return response, nil
}