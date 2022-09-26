package webapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"crypto/tls"

	// "log"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/jellyqwq/Paimon/requests"
	"github.com/jellyqwq/Paimon/tools"
	"github.com/jellyqwq/Paimon/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	response, err := requests.Bronya("GET", "https://fanyi.youdao.com/", nil, nil, "")
	if err != nil {
		return "", err
	}
	compile := regexp.MustCompile(`(?P<OUTFOX_SEARCH_USER_ID>OUTFOX_SEARCH_USER_ID=-?[0-9]+@[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+);`)
	// log.Println(response.Header["Set-Cookie"][0])
	dict := tools.GetParamsOneDimension(compile, response.Header["Set-Cookie"][0])

	t := tools.Md5(UA)
	un := time.Now().UnixNano()
	r := un / 1e6
	rand.Seed(un)
	i := strconv.FormatInt(r, 10) + strconv.FormatInt(int64(rand.Intn(10)), 10)

	ts := r
	bv := t
	salt := i
	sign := tools.Md5("fanyideskweb" + word + i + "Ygy_4c=r#e#4EX^NUGUc5")

	data := map[string]string {
		"i": word,
		"from": "AUTO",
		"to": "AUTO",
		"smartresult": "dict",
		"client": "fanyideskweb",
		"salt": salt,
		"sign": sign,
		"lts": strconv.FormatInt(int64(ts), 10),
		"bv": bv,
		"doctype": "json",
		"version": "2.1",
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
		"Cookie": fmt.Sprintf("%v; OUTFOX_SEARCH_USER_ID_NCOO=%v; ___rl__test__cookies=%v", dict["OUTFOX_SEARCH_USER_ID"], 2147483647*rand.Float64(), time.Now().UnixNano()/1e6),
	}

	response, err = requests.Bronya("POST", "https://fanyi.youdao.com/translate_o?smartresult=dict&smartresult=rule", headers, data, "")
	if err != nil {
		return "", err
	}

	body := response.Body
	var youdao YoudaoTranslation
	json.Unmarshal(body, &youdao)

	if youdao.ErrorCode != 0 {
		return "", fmt.Errorf("request params error, %v", string(body))
	}

	return youdao.TranslateResult[0][0].Tgt, nil
}

type Params struct {
    Bot *tgbotapi.BotAPI
    Conf config.Config
}

// Elysia music♪
func (params *Params) YoutubeSearch(query string) ([]interface{}, error) {
	data := fmt.Sprintf(`{
		"context": {
			"client": {
				"clientName": "WEB",
				"clientVersion": "2.20220617.00.00"
			}
		},
		"query": "%v"
	}`, query)
	headers := map[string]string{
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	}
	response, err := requests.Bronya("POST", "https://www.youtube.com/youtubei/v1/search", headers, nil, data)
	if err != nil {
		return nil, err
	}
	jsonRet := map[string]interface{}{}
	err = json.Unmarshal(response.Body, &jsonRet)
	if err != nil {
		return nil, err
	}

	contents := jsonRet["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"].(map[string]interface{})["sectionListRenderer"].(map[string]interface{})["contents"].([]interface{})[0].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"]

	var result []interface{}
	
	for _, c := range contents.([]interface{}) {
		if c.(map[string]interface{})["videoRenderer"] != nil {
			videoRenderer := c.(map[string]interface{})["videoRenderer"]
			videoId := videoRenderer.(map[string]interface{})["videoId"].(string)

			title := videoRenderer.(map[string]interface{})["title"].(map[string]interface{})["runs"].([]interface{})[0].(map[string]interface{})["text"].(string)
			// performer := videoRenderer.(map[string]interface{})["ownerText"].(map[string]interface{})["runs"].([]interface{})[0].(map[string]interface{})["text"].(string)
			result = append(result, (tgbotapi.NewInlineQueryResultAudio(
				videoId,
				// fmt.Sprintf("http://youtube.mp3.jellyqwq.com/ytb/%v", videoId),
				fmt.Sprintf("%vy2mate/%v", params.Conf.TelegramWebHook.Url, videoId), 
				title,
			)))
		}
	}
	return result, nil
}

func (params *Params) Y2mate(writer http.ResponseWriter, request *http.Request) {
	str := request.URL.String()
	videoId := str[8:]
	log.Println("videoId:", videoId)

	if videoId != "" {
		headers := map[string]string{
			"accept": "application/json, text/javascript, */*; q=0.01",
			"accept-encoding": "gzip, deflate, br",
			"accept-language": "zh-US,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-CN;q=0.6,ja-CN;q=0.5,ja;q=0.4",
			"cache-control": "no-cache",
			"pragma": "no-cache",
			"sec-ch-ua": "\"Google Chrome\";v=\"105\", \"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"105\"",
			"sec-ch-ua-mobile": "?0",
			"sec-ch-ua-platform": "Windows",
			"sec-fetch-dest": "empty",
			"sec-fetch-mode": "cors",
			"sec-fetch-site": "same-origin",
			"x-requested-with": "XMLHttpRequest",
			"origin": "https://y2mate.tools",
			"referer": "https://y2mate.tools/en57db",
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
		}
		data := map[string]string{
			"url": "https://www.youtube.com/watch?v=" + videoId,
			"q_auto": "0",
			"ajax": "1",
		}
		response , err := requests.Bronya("POST", "https://y2mate.tools/mates/en/analyze/ajax", headers, data, "")
		if err != nil {
			log.Println(err)
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
		// log.Println("result: ", result)

		y2mateCompile := regexp.MustCompile(`"https://converter\.quora-wiki\.com/#url=(?P<url>[0-9a-zA-Z=]+)".*?data-ftype="m4a"`)
		paramsMap := tools.GetParamsOneDimension(y2mateCompile, result)
		url := paramsMap["url"]
		if url == "" {
			log.Println("cannot find url")
			return
		}
		fmt.Println("y2mate", url)
		
		data = map[string]string{
			"url": url,
		}
		headers = map[string]string{
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
			"origin": "https://converter.quora-wiki.com",
			"referer": "https://converter.quora-wiki.com/",
		}
		response, err = requests.Bronya("POST", "https://converter.quora-wiki.com/", headers, data, "")
		if err != nil {
			fmt.Println(err)
		}
		audioJson := map[string]interface{}{}
		err = json.Unmarshal(response.Body, &audioJson)
		if err != nil {
			log.Println(err)
			return
		}
		audioUrl := audioJson["url"].(string)
		log.Println("audioUrl: ", audioUrl)

		// https://github.com/MeteorsLiu/FastPaimon/blob/c4704dc6f60cb9fa7181bcf6a73254e4e10cc835/httpServer/http.go#L98
		pr, pw := io.Pipe()
		//Async Fetch the audio
		cl := make(chan string)
		go func() {
			defer pw.Close()
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{
				Timeout: 	1 * time.Hour,
				Transport: 	tr,
			}
	
			req, err := http.NewRequest("GET", audioUrl, nil)
			if err != nil {
				log.Println(err)
				cl <- "0"
				return
			}
			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
			req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	
			resp, err := client.Do(req)
	
			if err != nil {
				log.Println(err)
				cl <- "0"
				return
			}
			defer resp.Body.Close()
			cl <- strconv.FormatInt(resp.ContentLength, 10)
	
			io.Copy(pw, resp.Body)
		}()
		_len := <-cl
		if _len == "0" {
			return
		}
		writer.Header().Set("Content-Length", _len)
		writer.Header().Set("Content-Type", "audio/mpeg")
		io.Copy(writer, pr)
	}
}
