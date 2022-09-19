package requests

import (
	"bytes"
	"net/http"
	nurl "net/url"
	"strconv"
)

func Post(url string, headers map[string]string, data map[string]interface{}) (*http.Response, error) {
	// data preprogressing, the value only support (string, int, float64) 
	var request *http.Request
	var err error
	if data != nil {
		dataForm := nurl.Values{}
		for key, val := range data {
			value := ""
			switch val := val.(type) {
				case string: value = string(val)
				case int: value = strconv.FormatInt(int64(val), 10)
				case int64: value = strconv.FormatInt(val, 10)
				case float64: value = strconv.FormatFloat(val, 'G', -1, 64)
			}
			dataForm.Set(key, value)
		}
	
		dataString := dataForm.Encode()
		dataByte := []byte(dataString)
		dataReader := bytes.NewReader(dataByte)
		request, err = http.NewRequest("POST", url, dataReader)
	} else {
		request, err = http.NewRequest("POST", url, nil)
	}
	if err != nil {
		return nil, err
	}
	
	for key, val := range headers {
		request.Header.Set(key, val)
	}

	client := &http.Client{}
	return client.Do(request)
}

func Get(url string, headers map[string]string, data map[string]interface{}) (*http.Response, error) {
	// data preprogressing, the value only support (string, int, float64) 
	var request *http.Request
	var err error
	if data != nil {
		dataForm := nurl.Values{}
		for key, val := range data {
			value := ""
			switch val := val.(type) {
				case string: value = string(val)
				case int: value = strconv.FormatInt(int64(val), 10)
				case int64: value = strconv.FormatInt(val, 10)
				case float64: value = strconv.FormatFloat(val, 'G', -1, 64)
			}
			dataForm.Set(key, value)
		}
	
		dataString := dataForm.Encode()
		dataByte := []byte(dataString)
		dataReader := bytes.NewReader(dataByte)
		request, err = http.NewRequest("GET", url, dataReader)
	} else {
		request, err = http.NewRequest("GET", url, nil)
	}
	if err != nil {
		return nil, err
	}
	
	for key, val := range headers {
		request.Header.Set(key, val)
	}

	client := &http.Client{}
	return client.Do(request)
}