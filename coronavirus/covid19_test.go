package coronavirus

import (
	"encoding/json"
	"fmt"
	"testing"

	"log"

	"github.com/jellyqwq/Paimon/requests"
)

func Test(t *testing.T) {
	res, err := requests.Bronya("GET", fmt.Sprintf("https://opendata.baidu.com/data/inner?resource_id=5653&query=%v%v新型肺炎最新动态&alr=1", "广东", "广州"), nil, nil, nil, false)
	if err != nil {
		log.Println(err)
	}
	js := map[string]interface{}{}
	err = json.Unmarshal(res.Body, &js);
	if err != nil {
		log.Println(err)
	}
	log.Println(string(res.Body))
	log.Println(js)
	log.Println(js["Result"])
	data_list := js["Result"].([]interface{})[0].(map[string]interface{})["DisplayData"].(map[string]interface{})["resultData"].(map[string]interface{})["tplData"].(map[string]interface{})["data_list"].([]interface{})
	log.Println(data_list)
	for _, tli := range data_list {
		li := tli.(map[string]interface{})
		log.Println(li["total_desc"].(string), li["total_num"].(string))
	}

}