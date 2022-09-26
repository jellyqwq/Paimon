package requests

import (
	// "bytes"
	"fmt"
	"io"
	"time"

	// "log"
	"net/http"
	nurl "net/url"
	"strconv"
	"strings"
)

type Response struct {
	StatusCode int
	Header     map[string][]string
	Body       []byte
	Contentlength int64
}

// Post and Get methods
func Bronya(method string, url string, headers map[string]string, data map[string]string, json string) (*Response, error) {
	var request *http.Request
	var err error
	if method == "POST" {
		if json != "" {
			request, err = http.NewRequest("POST", url, strings.NewReader(json))
			request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		} else if data != nil {
			formdata := nurl.Values{}
			for k, v := range data {
				formdata.Set(k,v)
			}
			dataString := formdata.Encode()
			// dataByte := []byte(dataString)
			contenLength := len(dataString)
			// log.Println(dataString)
			request, err = http.NewRequest("POST", url, strings.NewReader(dataString))
			request.Header.Set("Content-Length", strconv.FormatInt(int64(contenLength), 10))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		} else {
			return nil, fmt.Errorf("data and json only can choose one")
		}
	} else if method == "GET" {
		// GET请求一般不携带内容
		request, err = http.NewRequest("GET", url, nil)
	} else {
		return nil, fmt.Errorf("method must be POST or GET")
	}

	if err != nil {
		return nil, err
	}

	// set request headers
	for key, val := range headers {
		request.Header.Set(key, val)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	response := &Response{}
	response.StatusCode = res.StatusCode
	response.Header = res.Header
	response.Contentlength = res.ContentLength

	response.Body, err = io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	return response, nil
}
