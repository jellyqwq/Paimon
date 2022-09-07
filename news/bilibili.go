package news

import (
	"net/http"
	"io"
	"encoding/json"
	"strconv"
	"fmt"
)

// B站热搜
func BiliHotWords() (string, error) {
	request, err := http.NewRequest("GET", "http://s.search.bilibili.com/main/hotword", nil)
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

	result := ""

	if jsonRet["code"].(float64) == 0.0 {
		
		list := jsonRet["list"]
		
		for i, li := range list.([]interface{}) {
			
			limap := li.(map[string]interface{})
			
			show_name := limap["show_name"].(string)
			id := strconv.FormatFloat(limap["id"].(float64), 'f', 0, 64)

			show_name = string(show_name)
			if i != 9 {
				result += id + ".[" + show_name + "](https://search.bilibili.com/all?keyword=" + show_name + ")\n"
			} else {
				result += id + ".[" + show_name + "](https://search.bilibili.com/all?keyword=" + show_name + ")"
			}
		}
		return result, nil
	} else {
		return "", fmt.Errorf("GG")
	}
}