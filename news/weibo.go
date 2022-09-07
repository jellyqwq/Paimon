package news

import (
	"encoding/json"
	"fmt"
	// "strconv"

	"io"
	"net/http"
)

// 微博热搜有俩api
// https://weibo.com/ajax/side/hotSearch
// https://weibo.com/ajax/statuses/hot_band

func WeiboHotWords() (string, error) {
	request ,err := http.NewRequest("GET", "https://weibo.com/ajax/statuses/hot_band", nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	jsonRet := map[string]interface{}{}
	if err := json.Unmarshal(body, &jsonRet); err != nil {
		return "", err
	}

	if int64(jsonRet["http_code"].(float64)) == 200 {
		data := jsonRet["data"].(map[string]interface{})
		hotgov := data["hotgov"].(map[string]interface{})["word"].(string)
		band_list := data["band_list"]

		result := "Top.[" + hotgov + "](https://s.weibo.com/weibo?q=" + hotgov + ")\n"
		var index int64 = 1
		for _, li := range band_list.([]interface{}) {
			band := li.(map[string]interface{})
			if band["ad_type"] != nil || band["category"] == "综艺" {
				continue
			}
			// result += strconv.FormatInt(index, 10) + ".[" + band["word"].(string) + "](" + "https://s.weibo.com/weibo?q=" + band["word"].(string) + ")\n"
			result += fmt.Sprintf("%d.[%s](https://s.weibo.com/weibo?q=%s)\n", index, band["word"].(string), band["word"].(string))
			if index == 10 {
				break
			}
			index += 1
		}
		return result, err


	} else {
		return "", fmt.Errorf("weibo api error")
	}
	
}