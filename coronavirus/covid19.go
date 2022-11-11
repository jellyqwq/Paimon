package coronavirus

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	
	"github.com/jellyqwq/Paimon/requests"
	"github.com/jellyqwq/Paimon/tools"
)

// General
type VirusDetailed struct {
	New struct {
		// 新增确诊
		Diagnose struct {
			Abroad                int `json:"abroad"`
			Mainland              int `json:"mainland"`
			AbroadFromAsymptoma   int `json:"abroad_from_asymptoma"`
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
	Title string `json:"title"`
	Time  string `json:"time"`
	New   struct {
		// 新增确诊
		Diagnose struct {
			Total                 int `json:"total"`
			Abroad                int `json:"abroad"`
			Mainland              int `json:"mainland"`
			AbroadFromAsymptoma   int `json:"abroad_from_asymptoma"`
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
	ProvinceDetailed map[string]*VirusDetailed
}

var (
	CompilePageTitle     = regexp.MustCompile(`<div class="tit">(?P<title>.*?)<\/div>`)
	CompileFormatPage    = regexp.MustCompile(`<div class="con" id="xw_box">(?s:.*)<div class="fx fr">`)
	CompilePageLabel     = regexp.MustCompile(`>(.*?)<`)
	CompileProvinceCount = regexp.MustCompile(`(?P<Province>.*?)(?P<Count>\d+)例`)
	// 新增数据匹配
	CompileDiagenose = regexp.MustCompile(`新增确诊病例(?P<NewDiagenoseTotal>\d+)例.*?境外输入病例(?P<NewDiagnoseAbroadTotal>\d+)例（(?P<NewDiagnoseAbroadString>.*?)），含(?P<NewDiagnoseAbroadFromAsymptoma>\d+)例由无症状感染者转为确诊病例（(?P<NewDiagnoseAbroadFromAsymptomaString>.*?)）；本土病例(?P<NewDiagnoseMainlandTotal>\d+)例（(?P<NewDiagnoseMainlandString>.*?)），含(?P<NewDiagnoseMainlandFromAsymptoma>\d+)例由无症状感染者转为确诊病例（(?P<NewDiagnoseMainlandFromAsymptomaString>.*?)）。(?:(?P<NewDeaths>无新增死亡病例)。)(?:(?P<NewSuspected>无新增疑似病例)。)`)

	Core             KernelVirus
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
func GetAnnouncementList() (AnnouncementList *[]map[string]string, err error) {
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

	AnnouncementList = tools.GetParamsMultiDimension(compileCovid19, str)

	return AnnouncementList, nil
}

// 提取详情页文本
func GetAnnouncementString(url string) (string, error) {
	// request detailed page
	rep, err := requests.Bronya("GET", "http://www.nhc.gov.cn"+url, nil, nil, nil, false)
	if err != nil {
		return "", err
	}

	if rep.StatusCode != 200 {
		return "", fmt.Errorf("StatusCode: %v", rep.StatusCode)
	}

	temp := CompileFormatPage.FindString(string(rep.Body))

	temp2 := CompilePageLabel.FindAllStringSubmatch(temp, -1)

	var OriginalString string
	for _, t := range temp2 {
		if t[1] != "" {
			OriginalString += strings.ReplaceAll(t[1], "\n", "")
		}
	}
	return OriginalString, nil
}

// 获取今日详情页内容
func GetTodayAnnouncement() (string, error) {
	AnnouncementList, err := GetAnnouncementList()
	if err != nil {
		return "", err
	}

	url := (*AnnouncementList)[0]["url"]

	Core.Time = (*AnnouncementList)[0]["time"]
	Core.Title = (*AnnouncementList)[0]["title"]

	text, err := GetAnnouncementString(url)
	if err != nil {
		return "", err
	}

	return text, nil
}

// 新增写入
func NewDataWrite(OriginalString string) {
	// NewDeaths, NewSuspected
	ParamsNewData := tools.GetParamsOneDimension(CompileDiagenose, OriginalString)

	// 新增确诊病例
	if ParamsNewData["NewDiagenoseTotal"] != "" {
		num, err := strconv.Atoi(ParamsNewData["NewDiagenoseTotal"])
		if err != nil {
			log.Println(err)
		} else {
			Core.New.Diagnose.Total = num
		}
	}

	// 新增境外输入
	if ParamsNewData["NewDiagnoseAbroadTotal"] != "" {
		num, err := strconv.Atoi(ParamsNewData["NewDiagnoseAbroadTotal"])
		if err != nil {
			log.Println(err)
		} else {
			Core.New.Diagnose.Abroad = num
		}
	}

	// 新增境外输入中包含的无症状转确诊
	if ParamsNewData["NewDiagnoseAbroadFromAsymptoma"] != "" {
		num, err := strconv.Atoi(ParamsNewData["NewDiagnoseAbroadFromAsymptoma"])
		if err != nil {
			log.Println(err)
		} else {
			Core.New.Diagnose.AbroadFromAsymptoma = num
		}
	}

	// 新增本土病例
	if ParamsNewData["NewDiagnoseMainlandTotal"] != "" {
		num, err := strconv.Atoi(ParamsNewData["NewDiagnoseMainlandTotal"])
		if err != nil {
			log.Println(err)
		} else {
			Core.New.Diagnose.Mainland = num
		}
	}

	// 新增本土病例中包含的无症状转确诊
	if ParamsNewData["NewDiagnoseMainlandFromAsymptoma"] != "" {
		num, err := strconv.Atoi(ParamsNewData["NewDiagnoseMainlandFromAsymptoma"])
		if err != nil {
			log.Println(err)
		} else {
			Core.New.Diagnose.MainlandFromAsymptoma = num
		}
	}

	Core.ProvinceDetailed = make(map[string]*VirusDetailed)

	// 各省份新增境外输入病例
	if ParamsNewData["NewDiagnoseAbroadString"] != "" {
		str := ParamsNewData["NewDiagnoseAbroadString"]
		list := strings.Split(str, "，")

		for _, i := range list {
			dict := tools.GetParamsOneDimension(CompileProvinceCount, i)
			province := dict["Province"]
			Count, err := strconv.Atoi(dict["Count"])
			if err != nil {
				log.Println(err)
				continue
			}
			if Core.ProvinceDetailed[province] == nil {
				Core.ProvinceDetailed[province] = &VirusDetailed{}
			}

			TempStruct := Core.ProvinceDetailed[province]
			TempStruct.New.Diagnose.Abroad = Count
			Core.ProvinceDetailed[province] = TempStruct
		}
	}

	// 各省份新增境外输入病例中无症状转确诊
	if ParamsNewData["NewDiagnoseAbroadFromAsymptomaString"] != "" {
		str := ParamsNewData["NewDiagnoseAbroadFromAsymptomaString"]
		list := strings.Split(str, "，")
		for _, i := range list {
			dict := tools.GetParamsOneDimension(CompileProvinceCount, i)
			province := dict["Province"]
			Count, err := strconv.Atoi(dict["Count"])
			if err != nil {
				log.Println(err)
				continue
			}
			if Core.ProvinceDetailed[province] == nil {
				Core.ProvinceDetailed[province] = &VirusDetailed{}
			}

			TempStruct := Core.ProvinceDetailed[province]
			TempStruct.New.Diagnose.AbroadFromAsymptoma = Count
			Core.ProvinceDetailed[province] = TempStruct
		}
	}

	// 各省份新增本土病例
	if ParamsNewData["NewDiagnoseMainlandString"] != "" {
		str := ParamsNewData["NewDiagnoseMainlandString"]
		list := strings.Split(str, "，")

		for _, i := range list {
			dict := tools.GetParamsOneDimension(CompileProvinceCount, i)
			province := dict["Province"]
			Count, err := strconv.Atoi(dict["Count"])
			if err != nil {
				log.Println(err)
				continue
			}
			if Core.ProvinceDetailed[province] == nil {
				Core.ProvinceDetailed[province] = &VirusDetailed{}
			}

			TempStruct := Core.ProvinceDetailed[province]
			TempStruct.New.Diagnose.Mainland = Count
			Core.ProvinceDetailed[province] = TempStruct
		}
	}

	// 各省份新增本土病例中无症状转确诊
	if ParamsNewData["NewDiagnoseMainlandFromAsymptomaString"] != "" {
		str := ParamsNewData["NewDiagnoseMainlandFromAsymptomaString"]
		list := strings.Split(str, "，")
		for _, i := range list {
			dict := tools.GetParamsOneDimension(CompileProvinceCount, i)
			province := dict["Province"]
			Count, err := strconv.Atoi(dict["Count"])
			if err != nil {
				log.Println(err)
				continue
			}
			if Core.ProvinceDetailed[province] == nil {
				Core.ProvinceDetailed[province] = &VirusDetailed{}
			}

			TempStruct := Core.ProvinceDetailed[province]
			TempStruct.New.Diagnose.MainlandFromAsymptoma = Count
			Core.ProvinceDetailed[province] = TempStruct
		}
	}
}

func Entry() (*KernelVirus, error) {
	OriginalString, err := GetTodayAnnouncement()
	if err != nil {
		return nil, err
	}

	NewDataWrite(OriginalString)

	return &Core, nil
}
