package webapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/jellyqwq/Paimon/tools"
)

type YoudaoTranslation struct {
	ErrorCode       int `json:"errorCode"`
	TranslateResult [][]struct {
		Tgt string `json:"tgt"`
		Src string `json:"src"`
	} `json:"translateResult"`
	Type        string `json:"type"`
	SmartResult struct {
		Entries []string `json:"entries"`
		Type    int      `json:"type"`
	} `json:"smartResult"`
}


func RranslateByYouDao(word string) (string, error) {
	UA := "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"

	preClient := &http.Client{}
	request, err := http.NewRequest("GET", "https://fanyi.youdao.com/", nil)
	if err != nil {
		return "", err
	}
	response, err := preClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	compile := regexp.MustCompile(`(?P<OUTFOX_SEARCH_USER_ID>OUTFOX_SEARCH_USER_ID=-?[0-9]+@[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+);`)
	dict := tools.GetParamsOneDimension(compile, response.Header["Set-Cookie"][0])
	
	t := tools.Md5(UA)
	un := time.Now().UnixNano()
	r := un / 1e6
	rand.Seed(un)
	i := r + int64(rand.Intn(10))

	ts := r
	bv := t
	salt := i
	sign := tools.Md5("fanyideskweb" + word + strconv.FormatInt(i, 10) + "Ygy_4c=r#e#4EX^NUGUc5")

	form_data := make(map[string]string)
	form_data["i"] = word
	form_data["from"] = "AUTO"
	form_data["to"] = "AUTO"
	form_data["smartresult"] = "dict"
	form_data["client"] = "fanyideskweb"
	form_data["salt"] = strconv.FormatInt(salt, 10)
	form_data["sign"] = sign
	form_data["lts"] = strconv.FormatInt(ts, 10)
	form_data["bv"] = bv
	form_data["doctype"] = "json"
	form_data["version"] = "2.1"
	form_data["keyfrom"] = "fanyi.web"
	form_data["action"] = "FY_BY_REALTlME"

	formValues := url.Values{}
	for key, value := range form_data {
		formValues.Set(key, value)
	}

	formDataStr := formValues.Encode()
	formDataBytes := []byte(formDataStr)
	formBytesReader := bytes.NewReader(formDataBytes)

	url := "https://fanyi.youdao.com/translate_o?smartresult=dict&smartresult=rule"
	client := &http.Client{}
	request, err = http.NewRequest("POST", url, formBytesReader)

	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("User-Agent", UA)
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("Origin", "https://fanyi.youdao.com")
	request.Header.Set("Referer", "https://fanyi.youdao.com/")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Connection", "keep-alive")
	// request.Header.Set("Content-Length", strconv.FormatInt(int64(content_length), 10))
	request.Header.Set("Cookie", fmt.Sprintf("%v; OUTFOX_SEARCH_USER_ID_NCOO=%v; ___rl__test__cookies=%v", dict["OUTFOX_SEARCH_USER_ID"], 2147483647 * rand.Float64(), time.Now().UnixNano() / 1e6))
	request.Header.Set("Host", "fanyi.youdao.com")

	response, err = client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var youdao YoudaoTranslation
	json.Unmarshal(body, &youdao)

	if youdao.ErrorCode != 0 {
		return "", fmt.Errorf("request params error, %v", string(body))
	}

	return youdao.TranslateResult[0][0].Tgt, nil
}