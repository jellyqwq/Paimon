package coronavirus

import (
	"log"
	"strings"

	"testing"

	"github.com/jellyqwq/Paimon/requests"
	"github.com/jellyqwq/Paimon/tools"
)

func Test(t *testing.T) {
	url := "http://www.nhc.gov.cn/xcs/yqtb/202211/db29a8e77cc9449998ebf2c2548d096a.shtml"
	rep, err := requests.Bronya("GET", url, nil, nil, nil, false)
	if err != nil {
		log.Println(err)
	}

	if rep.StatusCode != 200 {
		log.Printf("StatusCode: %v", rep.StatusCode)
	}

	// log.Println(string(rep.Body))
	params := tools.GetParamsOneDimension(CompilePageTitle, string(rep.Body))
	log.Println(params["title"])
	m := CompileFormatPage.FindString(string(rep.Body))
	// log.Println(m)
	r := CompilePageLabel.FindAllStringSubmatch(m, -1)
	OriginalString := ""
	for _, s := range r {
		if s[1] != "" {
			OriginalString += strings.ReplaceAll(s[1], "\n", "")
		}
	}
	log.Println(OriginalString)

	
}
