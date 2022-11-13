package coronavirus

import (
	"log"
	"regexp"

	"testing"
	"github.com/jellyqwq/Paimon/requests"
	// "github.com/jellyqwq/Paimon/tools"
)

func Test(t *testing.T) {
	headers := map[string]string{
		"Accept":                    "text/html",
		// "Accept":                    "*/*",
		// "Accept-Encoding":           "deflate",
		"Accept-Language":           "zh-US,zh;q=0.9",
		"Cache-Control":             "no-cache",
		"Host":                      "www.nhc.gov.cn",
		"Pragma":                    "no-cache",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	}
	response, err := requests.Bronya("GET", "http://www.nhc.gov.cn/xcs/yqtb/list_gzbd.shtml", headers, nil, nil, false)
	if err != nil {
		log.Println(err)
	}

	str := string(response.Body)

	compileCovid19 := regexp.MustCompile(`<a href="(?P<url>.*?)".*?title='(?P<title>.*?)'.*?<span class="ml">(?P<time>.*?)</span>`)

	AnnouncementList := GetParamsMultiDimension(compileCovid19, str)
	log.Println((AnnouncementList))
}

func GetParamsMultiDimension(compile *regexp.Regexp, str string) (paramsMap *[]map[string]string) {
	match := compile.FindAllStringSubmatch(str, -1)

	var paramsList []map[string]string
	for _, m := range match {
		temp := map[string]string{}
		for i, name := range compile.SubexpNames() {
			if i > 0 && i <= len(m) {
				temp[name] = m[i]
			}
		}
		paramsList = append(paramsList, temp)
	}
	return &paramsList
}
