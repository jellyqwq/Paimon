package news

import (
	"log"
	"net/http"
	"io"
	"encoding/json"
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

	// log.Printf("%+v\n", body)

	jsonRet := map[string]interface{}{}
	if err := json.Unmarshal(body, &jsonRet); err != nil {
		return "", err
	}
	// log.Printf("%+v\n", jsonRet)

	
	result := ""

	if jsonRet["code"].(float64) == 0.0 {
		log.Printf("%+v\n", int64(jsonRet["timestamp"].(float64)))
		list := jsonRet["list"]
		// log.Printf("%+v\n", list)
		for i, li := range list.([]interface{}) {
			log.Printf("%v: %v\n", i, li)
			a := li.(map[string]interface{})
			
			show_name := a["show_name"].(string)
			id := fmt.Sprintf("%d", int64(a["id"].(float64)))

			if i != 9 {
				result = result + id + "." + string(show_name) + "\n"
			} else {
				result = result + id + "." + string(show_name)
			}
		}
		return result, nil
	} else {
		return "", fmt.Errorf("GG")
	}
}