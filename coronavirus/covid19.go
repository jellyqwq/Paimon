package coronavirus

import (
	"fmt"
	"regexp"

	// "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jellyqwq/Paimon/requests"
	"github.com/jellyqwq/Paimon/tools"
)


// General
type VirusDetailed struct {
	New struct {
		// 新增确诊
		Diagnose struct {
			Abroad                int `json:"abroad"`
			AbroadFromAsymptoma   int `json:"abroad_from_asymptoma"`
			Mainland              int `json:"mainland"`
			MainlandFromAsymptoma int `json:"mainland_from_asymptoma"`
		} `json:"diagnose"`
		// 新增死亡
		Deaths struct {
			Abroad   int `json:"abroad"`
			Mainland int `json:"mainland"`
		} `json:"deadth"`
		// 新增治愈
		Cure struct {
			Abroad   int `json:"abroad"`
			Mainland int `json:"mainland"`
		} `json:"cure"`
	} `json:"new"`
}

type KernelVirus struct {
	Title	string `json:"title"`
	New struct {
		// 新增确诊
		Diagnose struct {
			Total                 int `json:"total"`
			Abroad                int `json:"abroad"`
			AbroadFromAsymptoma   int `json:"abroad_from_asymptoma"`
			Mainland              int `json:"mainland"`
			MainlandFromAsymptoma int `json:"mainland_from_asymptoma"`
		} `json:"diagnose"`
		// 新增死亡
		Deaths struct {
			Total    int `json:"total"`
			Abroad   int `json:"abroad"`
			Mainland int `json:"mainland"`
		} `json:"deadth"`
		// 新增治愈
		Cure struct {
			Total    int `json:"total"`
			Abroad   int `json:"abroad"`
			Mainland int `json:"mainland"`
		} `json:"cure"`
	} `json:"new"`
	// 省份映射
	ProvinceDetailed map[string]VirusDetailed
}

var (
	CompilePageTitle  = regexp.MustCompile(`<div class="tit">(?P<title>.*?)<\/div>`)
	CompileFormatPage = regexp.MustCompile(`<div class="con" id="xw_box">(?s:.*)<div class="fx fr">`)
	CompilePageLabel = regexp.MustCompile(`>(.*?)<`)
	// 新增数据匹配
	CompileDiagenose = regexp.MustCompile(`新增确诊病例(?P<NewDiagenoseTotal>\d+)例.*?境外输入病例(?P<NewDiagnoseAbroadTotal>\d+)例（(?P<NewDiagnoseAbroadString>.*?)），含(?P<NewDiagnoseAbroadFromAsymptoma>\d+)例由无症状感染者转为确诊病例（(?P<NewDiagnoseAbroadFromAsymptomaString>.*?)）；本土病例(?P<NewDiagnoseMainlandTotal>\d+)例（(?P<NewDiagnoseMainlandString>.*?)），含(?P<NewDiagnoseMainlandFromAsymptoma>\d+)例由无症状感染者转为确诊病例（(?P<NewDiagnoseMainlandFromAsymptomaString>.*?)）。(?:(?P<NewDeaths>无新增死亡病例)。)(?:(?P<NewSuspected>无新增疑似病例)。)`)
	
	Core KernelVirus
	DetailedToSimple = map[string]string{
		"北京":  "京",
		"天津":  "津",
		"河北":  "冀",
		"山西":  "晋",
		"内蒙古": "蒙",
		"辽宁":  "辽",
		"吉林":  "吉",
		"黑龙江": "黑",
		"上海":  "沪",
		"江苏":  "苏",
		"浙江":  "浙",
		"安徽":  "皖",
		"福建":  "闽",
		"江西":  "赣",
		"山东":  "鲁",
		"河南":  "豫",
		"湖北":  "鄂",
		"湖南":  "湘",
		"广东":  "粤",
		"广西":  "桂",
		"海南":  "琼",
		"四川":  "川",
		"贵州":  "贵",
		"云南":  "云",
		"重庆":  "渝",
		"西藏":  "藏",
		"陕西":  "陕",
		"甘肃":  "甘",
		"青海":  "青",
		"宁夏":  "宁",
		"新疆":  "新",
		"香港":  "港",
		"澳门":  "澳",
		"台湾":  "台",
	}
)

// 请求数据显示的函数
//
// Three keys of map are url, title, time.
func GetCoronavirusList() (result *[]map[string]string, err error) {
	headers := map[string]string{
		"Accept":          "text/html",
		"Accept-Encoding": "deflate",
		"Accept-Language": "zh-US,zh;q=0.9",
		"Cache-Control":   "no-cache",
		"Host":            "www.nhc.gov.cn",
		"Pragma":          "no-cache",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	}
	response, err := requests.Bronya("GET", "http://www.nhc.gov.cn/xcs/yqtb/list_gzbd.shtml", headers, nil, nil, false)
	if err != nil {
		return nil, err
	}

	str := string(response.Body)
	compileCovid19 := regexp.MustCompile(`<a href="(?P<url>.*?)".*?title='(?P<title>.*?)'.*?<span class="ml">(?P<time>.*?)</span>`)

	result = tools.GetParamsMultiDimension(compileCovid19, str)

	return result, nil
}


// Parser
//
// Input a url that can request page and parse it to a province map if it is exist.
func ParseAnnouncement(url string) (map[string]VirusDetailed, error) {
	// request detailed page
	rep, err := requests.Bronya("GET", url, nil, nil, nil, false)
	if err != nil {
		return nil, err
	}

	if rep.StatusCode != 200 {
		return nil, fmt.Errorf("StatusCode: %v", rep.StatusCode)
	}

	return nil, fmt.Errorf("None")
}
