package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jellyqwq/Paimon/requests"
)

func main() {
	// block one
	query := "moon hola"
	data := fmt.Sprintf(`{
		"context": {
			"client": {
				"clientName": "WEB",
				"clientVersion": "2.20220617.00.00"
			}
		},
		"query": %s,
	}`, query)
	response, err := requests.Bronya("POST", "https://www.youtube.com/youtubei/v1/search", nil, data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(response.StatusCode)
	log.Println(response.Header)
	jsonRet := map[string]interface{}{}
	json.Unmarshal(response.Body, &jsonRet)
	log.Println(jsonRet)
	contents := jsonRet["contents"]
	// .(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"].(map[string]interface{})["sectionListRenderer"].(map[string][]interface{})["contents"]
	// [0].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"]
	log.Println(contents)



	// block two

}

