package webapi

import (
	"encoding/json"
	// "bytes"
	"fmt"
	"log"
	// "io"
	// "net/http"
	// "net/url"

	"regexp"
	"testing"

	"github.com/jellyqwq/Paimon/requests"
	"github.com/jellyqwq/Paimon/tools"
)

func Test(t *testing.T) {
	videoId := "xREK6gZxYLQ"
	data := map[string]string{
		"url": "https://www.youtube.com/watch?v=" + videoId,
		"q_auto": "0",
		"ajax": "1",
	}

	response, err := requests.Bronya("POST", "https://y2mate.tools/mates/en/analyze/ajax", nil, data, "")
	if err != nil {
		fmt.Println(err)
	}
	log.Println("StatusCode: ", response.StatusCode)
	jsonRet := map[string]interface{}{}
	err = json.Unmarshal(response.Body, &jsonRet)
	if err != nil {
		log.Println("ERROR: 获取url失败")
		log.Println("Response Body:", string(response.Body))
		return
	}
	result := jsonRet["result"].(string)
	log.Println("result: ", result)
	y2mateCompile := regexp.MustCompile(`"https://converter\.quora-wiki\.com/#url=(?P<url>[0-9a-zA-Z=]+)".*?data-ftype="m4a"`)
	paramsMap := tools.GetParamsOneDimension(y2mateCompile, result)
	url := paramsMap["url"]
	if url == "" {
		log.Println("cannot find url")
		return
	}
	fmt.Println("y2mate", url)
}