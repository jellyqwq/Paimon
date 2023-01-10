package webapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"log"
	"net/http"

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

	if len(youdao.TranslateResult) == 0 {
		return "", fmt.Errorf("translate result is empty")
	}
	if len(youdao.TranslateResult[0]) == 0 {
		return "", fmt.Errorf("translate result index[0] is empty")
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

// 米游社随机cos
type CosApi struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			Post struct {
				GameID     int      `json:"game_id"`
				PostID     string   `json:"post_id"`
				FForumID   int      `json:"f_forum_id"`
				UID        string   `json:"uid"`
				Subject    string   `json:"subject"`
				Content    string   `json:"content"`
				Cover      string   `json:"cover"`
				ViewType   int      `json:"view_type"`
				CreatedAt  int      `json:"created_at"`
				Images     []string `json:"images"`
				PostStatus struct {
					IsTop      bool `json:"is_top"`
					IsGood     bool `json:"is_good"`
					IsOfficial bool `json:"is_official"`
				} `json:"post_status"`
				TopicIds               []int         `json:"topic_ids"`
				ViewStatus             int           `json:"view_status"`
				MaxFloor               int           `json:"max_floor"`
				IsOriginal             int           `json:"is_original"`
				RepublishAuthorization int           `json:"republish_authorization"`
				ReplyTime              string        `json:"reply_time"`
				IsDeleted              int           `json:"is_deleted"`
				IsInteractive          bool          `json:"is_interactive"`
				StructuredContent      string        `json:"structured_content"`
				StructuredContentRows  []interface{} `json:"structured_content_rows"`
				ReviewID               int           `json:"review_id"`
				IsProfit               bool          `json:"is_profit"`
				IsInProfit             bool          `json:"is_in_profit"`
				UpdatedAt              int           `json:"updated_at"`
				DeletedAt              int           `json:"deleted_at"`
				PrePubStatus           int           `json:"pre_pub_status"`
				CateID                 int           `json:"cate_id"`
				ProfitPostStatus       int           `json:"profit_post_status"`
			} `json:"post"`
			Forum struct {
				ID        int         `json:"id"`
				Name      string      `json:"name"`
				Icon      string      `json:"icon"`
				GameID    int         `json:"game_id"`
				ForumCate interface{} `json:"forum_cate"`
			} `json:"forum"`
			Topics []struct {
				ID            int    `json:"id"`
				Name          string `json:"name"`
				Cover         string `json:"cover"`
				IsTop         bool   `json:"is_top"`
				IsGood        bool   `json:"is_good"`
				IsInteractive bool   `json:"is_interactive"`
				GameID        int    `json:"game_id"`
				ContentType   int    `json:"content_type"`
			} `json:"topics"`
			User struct {
				UID           string `json:"uid"`
				Nickname      string `json:"nickname"`
				Introduce     string `json:"introduce"`
				Avatar        string `json:"avatar"`
				Gender        int    `json:"gender"`
				Certification struct {
					Type  int    `json:"type"`
					Label string `json:"label"`
				} `json:"certification"`
				LevelExp struct {
					Level int `json:"level"`
					Exp   int `json:"exp"`
				} `json:"level_exp"`
				IsFollowing bool   `json:"is_following"`
				IsFollowed  bool   `json:"is_followed"`
				AvatarURL   string `json:"avatar_url"`
				Pendant     string `json:"pendant"`
			} `json:"user"`
			SelfOperation struct {
				Attitude    int  `json:"attitude"`
				IsCollected bool `json:"is_collected"`
			} `json:"self_operation"`
			Stat struct {
				ViewNum     int `json:"view_num"`
				ReplyNum    int `json:"reply_num"`
				LikeNum     int `json:"like_num"`
				BookmarkNum int `json:"bookmark_num"`
				ForwardNum  int `json:"forward_num"`
			} `json:"stat"`
			HelpSys struct {
				TopUp     interface{}   `json:"top_up"`
				TopN      []interface{} `json:"top_n"`
				AnswerNum int           `json:"answer_num"`
			} `json:"help_sys"`
			Cover struct {
				URL            string      `json:"url"`
				Height         int         `json:"height"`
				Width          int         `json:"width"`
				Format         string      `json:"format"`
				Size           string      `json:"size"`
				Crop           interface{} `json:"crop"`
				IsUserSetCover bool        `json:"is_user_set_cover"`
				ImageID        string      `json:"image_id"`
				EntityType     string      `json:"entity_type"`
				EntityID       string      `json:"entity_id"`
			} `json:"cover"`
			ImageList []struct {
				URL            string      `json:"url"`
				Height         int         `json:"height"`
				Width          int         `json:"width"`
				Format         string      `json:"format"`
				Size           string      `json:"size"`
				Crop           interface{} `json:"crop"`
				IsUserSetCover bool        `json:"is_user_set_cover"`
				ImageID        string      `json:"image_id"`
				EntityType     string      `json:"entity_type"`
				EntityID       string      `json:"entity_id"`
			} `json:"image_list"`
			IsOfficialMaster bool          `json:"is_official_master"`
			IsUserMaster     bool          `json:"is_user_master"`
			HotReplyExist    bool          `json:"hot_reply_exist"`
			VoteCount        int           `json:"vote_count"`
			LastModifyTime   int           `json:"last_modify_time"`
			RecommendType    string        `json:"recommend_type"`
			Collection       interface{}   `json:"collection"`
			VodList          []interface{} `json:"vod_list"`
			IsBlockOn        bool          `json:"is_block_on"`
			ForumRankInfo    interface{}   `json:"forum_rank_info"`
			LinkCardList     []interface{} `json:"link_card_list"`
		} `json:"list"`
		LastID   string `json:"last_id"`
		IsLast   bool   `json:"is_last"`
		IsOrigin bool   `json:"is_origin"`
	} `json:"data"`
}

func HoyoBBS() ([]string, error) {
	headers := map[string]string{
		"Accept": "application/json",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
		"Host": "bbs-api.mihoyo.com",
	}
	res, err := requests.Bronya("GET", "https://bbs-api.mihoyo.com/post/wapi/getForumPostList?forum_id=49&gids=2&page_size=20&sort_type=1", headers, nil, nil, false)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %v", res.StatusCode)
	}
	var temp CosApi
	json.Unmarshal(res.Body, &temp)
	ccore := [][]string{}
	for _, i := range temp.Data.List {
		core := []string{}
		for _, j := range i.ImageList {
			core = append(core, j.URL)
		}
		ccore = append(ccore, core)
	}
	rand.Seed(time.Now().Unix())
	n := rand.Intn(len(ccore))
	return ccore[n], nil
}

func MihoyoLiveCode() string {
	NewsListUrl := "https://bbs-api.miyoushe.com/post/wapi/getNewsList?gids=2&type=3"
	headers := map[string]string {
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		"Host": "bbs-api.miyoushe.com",
		"Accept": "application/json",
	}
	response, err := requests.Bronya("GET", NewsListUrl, headers, nil, nil, false)
	if err != nil {
		return "request NewsList error"
	}

	s := string(response.Body)
	compile := regexp.MustCompile(`《原神》(?P<version>[\.\d]+)版本前瞻.*?event-ys-live/index\.html\?act_id=(?P<act_id>[\dys]+)`)
	actIdMap := tools.GetParamsOneDimension(compile, s)
	if len(actIdMap) == 0 {
		return "no match act_id"
	}

	version := actIdMap["version"]
	act_id := actIdMap["act_id"]

	// 兑换码接口
	codeInfoUrl := "https://webstatic.mihoyo.com/bbslive/code/" + act_id + ".json"
	headers = map[string]string {
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		"Accept": "application/json",
	}
	response, err = requests.Bronya("GET", codeInfoUrl, headers, nil, nil, false)
	if response.StatusCode == 404 {
		return "live code is not release"
	}
	if err != nil {
		return "request code error"
	}
	
	jsonRes := []interface{}{}
	error := json.Unmarshal(response.Body, &jsonRes)
	if error != nil {
		return "json parse error"
	}

	result := fmt.Sprintf("《原神》[%v版本前瞻](%v)兑换码\n", version, "https://webstatic.mihoyo.com/bbs/event/event-ys-live/index.html?act_id="+ act_id)
	for _ , item := range jsonRes {
		imap := item.(map[string]interface{})
		result += "`" + imap["code"].(string) + "`" + MihoyoLiveCodeStringFormat(imap["title"].(string)) + "\n"
	}
	return result
}

func MihoyoLiveCodeStringFormat(s string) string {
	compileElementTag := regexp.MustCompile(`<.*?>`)
	s = compileElementTag.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "*", "\\*")
	return s
}