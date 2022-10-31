package main

import (
	// "encoding/json"
	// "bytes"
	// "fmt"
	// "time"
	// "log"
	// "strconv"

	// "fmt"
	// "bytes"
	"log"
	// "io"
	// "net/http"
	// "net/url"

	// "regexp"
	// "strings"
	// "testing"

	"github.com/jellyqwq/Paimon/requests"
	// "github.com/jellyqwq/Paimon/tools"
)

func main() {
	headers := map[string]string{
		"Accept": "text/html",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-US,zh;q=0.9",
		"Cache-Control": "no-cache",
		"Host": "www.nhc.gov.cn",
		"Pragma": "no-cache",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	}
	repsonse, err := requests.Bronya("GET", "http://www.nhc.gov.cn/xcs/yqtb/list_gzbd.shtml", headers, nil, "", false)
	if err != nil {
		log.Println(err)
		log.Println("error")
		return
	}

	log.Println(string(repsonse.Body))
	log.Println(repsonse.Header)
	log.Println("中文")
}