package webapi

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"

	// "net/http"
	// "net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/jellyqwq/Paimon/requests"
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

	response, err := requests.Get("https://fanyi.youdao.com/", nil, nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	compile := regexp.MustCompile(`(?P<OUTFOX_SEARCH_USER_ID>OUTFOX_SEARCH_USER_ID=-?[0-9]+@[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+);`)
	log.Println(response.Header["Set-Cookie"][0])
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

	data := map[string]interface{}{
		"i": word,
		"from": "AUTO",
		"to": "AUTO",
		"smartresult": "dict",
		"client": "fanyideskweb",
		"salt": salt,
		"sign": sign,
		"lts": ts,
		"bv": bv,
		"doctype": "json",
		"version": 2.1,
		"keyfrom": "fanyi.web",
		"action": "FY_BY_REALTlME",
	}

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"User-Agent": UA,
		"X-Requested-With": "XMLHttpRequest",
		"Origin": "https://fanyi.youdao.com",
		"Referer": "https://fanyi.youdao.com/",
		"Accept": "application/json",
		"Host": "fanyi.youdao.com",
		"Cookie": fmt.Sprintf("%v; OUTFOX_SEARCH_USER_ID_NCOO=%v; ___rl__test__cookies=%v", dict["OUTFOX_SEARCH_USER_ID"], 2147483647 * rand.Float64(), time.Now().UnixNano() / 1e6),
	}
	response, err = requests.Post("https://fanyi.youdao.com/translate_o?smartresult=dict&smartresult=rule", headers, data)
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

