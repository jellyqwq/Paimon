package webapi

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jellyqwq/Paimon/requests"
	"github.com/jellyqwq/Paimon/tools"
)

func MihoyoLiveCode() string {
	result := ""
	var (
		origin  string
		version string
		act_id  string
		compileGenshinLive = regexp.MustCompile(`《原神》(?P<version>[\.\d]+)版本前瞻.*?(?P<origin>https://webstatic\.mihoyo\.com.*?live.*?act_id=(?P<act_id>[\d\w]+))`)
		compileGenshinCodeVer = regexp.MustCompile(`"code_ver":"(?P<code_ver>.*?)"`)
	)
	// 获取公告资讯
	for last_id := 0; last_id < 520; last_id += 20 {
		NewsListUrl := fmt.Sprintf("https://bbs-api.miyoushe.com/post/wapi/getNewsList?gids=2&page_size=20&type=3&last_id=%v", last_id)
		headers := map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"Host":       "bbs-api.miyoushe.com",
			"Accept":     "application/json",
		}
		response, err := requests.Bronya("GET", NewsListUrl, headers, nil, nil, false)
		if err != nil {
			result = "request NewsList error"
			break
		}

		s := string(response.Body)
		actIdMap := tools.GetParamsOneDimension(compileGenshinLive, s)
		if len(actIdMap) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		origin = actIdMap["origin"]
		version = actIdMap["version"]
		act_id = actIdMap["act_id"]
		if act_id != "" {
			break
		}
	}
	if act_id == "" {
		return "cannot match act_id"
	}
	// fmt.Println("woc", origin, version, act_id)

	// 先获取码, 再获取
	// https://api-takumi.mihoyo.com/event/miyolive/index
	miyolive_index := "https://api-takumi.mihoyo.com/event/miyolive/index"
	headers := map[string]string{
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		"Accept":       "application/json",
		"Origin":       "https://webstatic.mihoyo.com",
		"Referer":      "https://webstatic.mihoyo.com/",
		"X-Rpc-Act_id": act_id,
	}
	response, err := requests.Bronya("GET", miyolive_index, headers, nil, nil, false)
	if response.StatusCode == 404 {
		if response.StatusCode == 404 {
			return "request miyolive.index error"
		}
		if err != nil {
			return fmt.Sprint(err)
		}
	}

	CodeVerMap := tools.GetParamsOneDimension(compileGenshinCodeVer, string(response.Body))
	if len(CodeVerMap) == 0 {
		return "CodeVerMap is empty"
	}
	code_ver := CodeVerMap["code_ver"]
	// fmt.Println("code_ver: ", code_ver)

	refreshCode := fmt.Sprintf("https://api-takumi-static.mihoyo.com/event/miyolive/refreshCode?version=%s&time=%d", code_ver, time.Now().Local().Unix())
	headers = map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		"Accept":     "application/json",
		"Origin":       "https://webstatic.mihoyo.com",
		"Referer":      "https://webstatic.mihoyo.com/",
		"X-Rpc-Act_id": act_id,
	}
	response, err = requests.Bronya("GET", refreshCode, headers, nil, nil, false)
	if response.StatusCode == 404 {
		return "live code is not release"
	}
	if err != nil {
		return fmt.Sprint(err)
	}
	
	// fmt.Println("ok", string(response.Body))

	jsonRes := map[string]interface{}{}
	err = json.Unmarshal(response.Body, &jsonRes)
	if err != nil {
		return "refresh json parse error"
	}

	if jsonRes["retcode"].(float64) != 0 {
		return jsonRes["message"].(string)
	}
	
	data := jsonRes["data"].(map[string]interface{})
	if data == nil {
		return "data is null"
	}

	code_list := data["code_list"].([]interface{})

	if len(code_list) == 0 {
		return fmt.Sprintf("《原神》[%v版本前瞻](%v)兑换码暂未生成🤔", version, origin)
	}

	result = fmt.Sprintf("《原神》[%v版本前瞻](%v)兑换码\n", version, origin)

	for _, item := range code_list {
		i := item.(map[string]interface{})
		result += "`" + i["code"].(string) + "` " + MihoyoLiveCodeStringFormat(i["title"].(string)) + "\n"
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
