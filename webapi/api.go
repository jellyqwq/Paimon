package webapi

import (
	"encoding/json"
	"fmt"

	// "io"
	"log"
	"net/http"

	// "crypto/tls"

	// "log"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/jellyqwq/Paimon/config"
	"github.com/jellyqwq/Paimon/requests"
	"github.com/jellyqwq/Paimon/tools"

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

	response, err := requests.Bronya("GET", "https://fanyi.youdao.com/", nil, nil, nil, false)
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

	data := map[string]string{
		"i":           word,
		"from":        "AUTO",
		"to":          "AUTO",
		"smartresult": "dict",
		"client":      "fanyideskweb",
		"salt":        salt,
		"sign":        sign,
		"lts":         strconv.FormatInt(int64(ts), 10),
		"bv":          bv,
		"doctype":     "json",
		"version":     "2.1",
		"keyfrom":     "fanyi.web",
		"action":      "FY_BY_REALTlME",
	}

	headers := map[string]string{
		"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
		"User-Agent":       UA,
		"X-Requested-With": "XMLHttpRequest",
		"Origin":           "https://fanyi.youdao.com",
		"Referer":          "https://fanyi.youdao.com/",
		"Accept":           "application/json",
		"Host":             "fanyi.youdao.com",
		"Cookie":           fmt.Sprintf("%v; OUTFOX_SEARCH_USER_ID_NCOO=%v; ___rl__test__cookies=%v", dict["OUTFOX_SEARCH_USER_ID"], 2147483647*rand.Float64(), time.Now().UnixNano()/1e6),
	}

	response, err = requests.Bronya("POST", "https://fanyi.youdao.com/translate_o?smartresult=dict&smartresult=rule", headers, data, nil, false)
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
	Bot  *tgbotapi.BotAPI
	Conf config.Config
}

// Elysia music♪
func (params *Params) YoutubeSearch(query string, inlineType string) ([]interface{}, error) {
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
	response, err := requests.Bronya("POST", "https://www.youtube.com/youtubei/v1/search", headers, nil, &data, false)
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

	timeStamp16 := time.Now().UnixNano() / 1e6
	timeStamp16String := strconv.FormatInt(timeStamp16, 10)
	for _, c := range contents.([]interface{}) {
		if c.(map[string]interface{})["videoRenderer"] != nil {
			videoRenderer := c.(map[string]interface{})["videoRenderer"]
			videoId := videoRenderer.(map[string]interface{})["videoId"].(string)

			title := videoRenderer.(map[string]interface{})["title"].(map[string]interface{})["runs"].([]interface{})[0].(map[string]interface{})["text"].(string)
			performer := videoRenderer.(map[string]interface{})["ownerText"].(map[string]interface{})["runs"].([]interface{})[0].(map[string]interface{})["text"].(string)

			lengthText := videoRenderer.(map[string]interface{})["lengthText"]
			var videoLength string
			var simpleText string
			if lengthText != nil {
				videoLength = lengthText.(map[string]interface{})["accessibility"].(map[string]interface{})["accessibilityData"].(map[string]interface{})["label"].(string)

				simpleText = lengthText.(map[string]interface{})["simpleText"].(string)
			} else {
				continue
			}

			// m1是y2mate.tools m2是www.y2mate.com
			aurl := ""
			if inlineType == "m1" {
				aurl = fmt.Sprintf("%vy2mate/tools/%v", params.Conf.TelegramWebHook.Url, videoId)
			} else if inlineType == "m2" {
				aurl = fmt.Sprintf("%vy2mate/com/%v", params.Conf.TelegramWebHook.Url, videoId)
			}

			audio := tgbotapi.NewInlineQueryResultAudio(
				timeStamp16String+"_"+videoId,
				aurl,
				title,
			)
			duration, err := countAudioSeconds(videoLength)
			if err != nil {
				log.Println(err)
			}
			audio.Performer = "[" + simpleText + "]" + performer
			audio.Duration = duration
			result = append(result, (audio))
		}
	}
	return result, nil
}

// count audio seconds
func countAudioSeconds(str string) (int, error) {
	compileHMS := regexp.MustCompile(`(?:(?P<H>[0-9]+) hours?, )?(?:(?P<M>[0-9]+) minutes?, )?(?P<S>[0-9]+) seconds?`)
	paramsMap := tools.GetParamsOneDimension(compileHMS, str)

	seconds, err := strconv.Atoi(paramsMap["S"])
	if err != nil {
		log.Println(err)
		return 0, err
	}

	if paramsMap["H"] != "" {
		h, err := strconv.Atoi(paramsMap["H"])
		if err != nil {
			log.Println(err)
			return 0, err
		}
		seconds += 3600 * h
	}

	if paramsMap["M"] != "" {
		m, err := strconv.Atoi(paramsMap["M"])
		if err != nil {
			log.Println(err)
			return 0, err
		}
		seconds += 60 * m
	}
	return seconds, nil
}

// y2mate.tools 处理视频转音频的函数
func (params *Params) Y2mateByTools(writer http.ResponseWriter, request *http.Request) {
	str := request.URL.String()
	videoId := str[14:]
	log.Println("videoId:", videoId)

	if videoId != "" {
		data := map[string]string{
			"url":    "https://www.youtube.com/watch?v=" + videoId,
			"q_auto": "0",
			"ajax":   "1",
		}
		response, err := requests.Bronya("POST", "https://y2mate.tools/mates/en/analyze/ajax", nil, data, nil, false)
		if err != nil {
			log.Println(err)
		}
		log.Println("StatusCode:", response.StatusCode)

		jsonRet := map[string]interface{}{}
		err = json.Unmarshal(response.Body, &jsonRet)
		if err != nil {
			log.Println("ERROR: y2mate.tools获取url失败", err)
			log.Println("Body:", string(response.Body))
			return
		}
		result := jsonRet["result"].(string)

		y2mateCompile := regexp.MustCompile(`"https://converter\.quora-wiki\.com/#url=(?P<url>.*?)".*?data-ftype="m4a"`)
		paramsMap := tools.GetParamsOneDimension(y2mateCompile, result)
		url := paramsMap["url"]
		if url == "" {
			log.Println("result:", result)
			log.Println("cannot find url in Y2mateByTools")
			return
		}

		data = map[string]string{
			"url": url,
		}
		headers := map[string]string{
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
			"origin":     "https://converter.quora-wiki.com",
			"referer":    "https://converter.quora-wiki.com/",
		}
		response, err = requests.Bronya("POST", "https://converter.quora-wiki.com/", headers, data, nil, false)
		if err != nil {
			log.Println(err)
			return
		}
		audioJson := map[string]interface{}{}
		err = json.Unmarshal(response.Body, &audioJson)
		if err != nil {
			log.Println(err)
			return
		}
		audioUrl := audioJson["url"].(string)
		log.Println("audioUrl: ", audioUrl)

		headers = map[string]string{
			"User-Agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
			"Connection":     "keep-alive",
			"Accept":         "*/*",
			"Range":          "bytes=0-",
			"Referer":        "https://converter.quora-wiki.com/",
			"Sec-Fetch-Site": "cross-site",
		}
		res, err := requests.Bronya("GET", audioUrl, headers, nil, nil, false)
		log.Println(res.StatusCode)
		log.Println(res.Header)
		log.Println(res.Contentlength)
		if err != nil {
			log.Println(err)
			return
		}

		writer.Header().Set("Content-Type", "audio/mpeg")
		writer.Header().Set("Content-Length", strconv.FormatInt(res.Contentlength, 10))
		writer.Header().Set("Keep-Alive", "timeout=15")
		writer.Write(res.Body)
		log.Println("OK", videoId)
	}
}

// y2mate.com
func (params *Params) Y2mateByCom(writer http.ResponseWriter, request *http.Request) {
	str := request.URL.String()
	videoId := str[12:]
	log.Println("videoId:", videoId)

	if videoId == "" {
		log.Println("videoId is empty")
		return
	}

	data := map[string]string{
		"url":    "https://www.youtube.com/watch?v=" + videoId,
		"q_auto": "0",
		"ajax":   "1",
	}
	response, err := requests.Bronya("POST", "https://www.y2mate.com/mates/en402/analyze/ajax", nil, data, nil, false)
	if err != nil {
		log.Println(err)
	}
	log.Println("StatusCode:", response.StatusCode)

	jsonRet := map[string]interface{}{}
	err = json.Unmarshal(response.Body, &jsonRet)
	if err != nil {
		log.Println("ERROR: y2mate.com获取url失败", err)
		log.Println("Body:", string(response.Body))
		return
	}
	result := jsonRet["result"].(string)

	y2mateCompile := regexp.MustCompile(`k__id.*?"(?P<id>.*?)"`)
	paramsMap := tools.GetParamsOneDimension(y2mateCompile, result)
	id := paramsMap["id"]
	if id == "" {
		log.Println("result:", result)
		log.Println("cannot find id")
		return
	}

	data = map[string]string{
		"type":     "youtube",
		"_id":      id,
		"v_id":     videoId,
		"ajax":     "1",
		"token":    "",
		"ftype":    "mp3",
		"fquality": "128",
	}

	response, err = requests.Bronya("POST", "https://www.y2mate.com/mates/convert", nil, data, nil, false)
	if err != nil {
		log.Println(err)
		return
	}
	audioJson := map[string]interface{}{}
	err = json.Unmarshal(response.Body, &audioJson)
	if err != nil {
		log.Println(err)
		return
	}
	audioResultString := audioJson["result"].(string)

	y2mateCompileComUrl := regexp.MustCompile(`href="(?P<url>.*?)"`)
	paramsMap = tools.GetParamsOneDimension(y2mateCompileComUrl, audioResultString)
	audioUrl := paramsMap["url"]

	if id == "" {
		log.Println("audioResultString:", audioResultString)
		log.Println("cannot find url in Y2mateByCom")
		return
	}

	headers := map[string]string{
		"User-Agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
		"Accept":         "*/*",
	}

	res, err := requests.Bronya("GET", audioUrl, headers, nil, nil, true)
	if err != nil {
		log.Println("audio url:", audioUrl)
		log.Println(err)
		return
	}
	log.Println(res.StatusCode)
	log.Println(res.Header)
	log.Println(res.Contentlength)

	writer.Header().Set("Content-Type", "audio/mpeg")
	writer.Header().Set("Content-Length", strconv.FormatInt(res.Contentlength, 10))
	writer.Header().Set("Keep-Alive", "timeout=15")
	writer.Write(res.Body)
	log.Println("OK", videoId)
}

// 汇率请求
func Finance(transType string) (string, error) {
	headers := map[string]string{
		"origin": "https://www.google.com",
		"pragma": "no-cache",
		"referer": "https://www.google.com/",
		"sec-ch-ua": "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
	}
	rep, err := requests.Bronya("GET", "https://www.google.com/finance/quote/" + transType, headers, nil, nil, false)
	if err != nil {
		return "", err
	}

	compileRT := regexp.MustCompile(`<div class=\"YMlKec fxKbKc\">(?P<rate>[\.0-9]+)</div>.*?class=\"ygUjEc\" jsname=\"Vebqub\">(?P<time>.*?) &middot;`)
	paramsMap := tools.GetParamsOneDimension(compileRT, string(rep.Body))
	rate := paramsMap["rate"]
	time := paramsMap["time"]
	
	content := fmt.Sprintf("%s\n[%s](https://www.google.com/finance/quote/%s) => %s", time, transType, transType, rate)
	return content, nil
}