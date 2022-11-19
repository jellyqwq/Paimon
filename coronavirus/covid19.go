package coronavirus

import (
	"encoding/json"
	"fmt"

	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jellyqwq/Paimon/config"
	"github.com/jellyqwq/Paimon/olog"
	"github.com/jellyqwq/Paimon/requests"
)

type Paimon struct {
	Log  *olog.Olog
	Conf *config.Config
	Bot  *tgbotapi.BotAPI
}

func GetData() (*Source, error) {
	res, err := requests.Bronya("GET", "https://api.inews.qq.com/newsqa/v1/query/inner/publish/modules/list?modules=localCityNCOVDataList,diseaseh5Shelf", nil, nil, nil, false)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %v", res.StatusCode)
	}

	var tempStruct Source
	err = json.Unmarshal(res.Body, &tempStruct)
	if err != nil {
		return nil, err
	}

	return &tempStruct, nil
}

// 国内疫情重写
func (core *Core) ChinaWin(source *Source) {
	ChinaTotal := source.Data.Diseaseh5Shelf.ChinaTotal
	core.China.UpdateTimeVirus = ChinaTotal.Mtime
	core.China.AccComfirm = ChinaTotal.Confirm
	core.China.NowComfirm = ChinaTotal.NowConfirm
	core.China.AddComfirm = ChinaTotal.ConfirmAdd
	core.China.NowComfirmLocal = ChinaTotal.LocalConfirm
	core.China.AddComfirmLocal = ChinaTotal.LocalConfirmAdd
	core.China.SevereCase = ChinaTotal.NowSevere
	core.China.AccCure = ChinaTotal.Heal
	core.China.AccDead = ChinaTotal.Dead
	core.China.AddDead = ChinaTotal.DeadAdd
	core.China.NowLocalAsymptoma = ChinaTotal.NowLocalWzz
	core.China.AddLocalAsymptoma = ChinaTotal.LocalWzzAdd
	core.China.UpdateTimeRisk = ChinaTotal.MRiskTime
	core.China.HighRiskAreaNum = ChinaTotal.HighRiskAreaNum
	core.China.MediumRiskAreaNum = ChinaTotal.MediumRiskAreaNum
}

// 省份解析
func (core *Core) ParseProvince(source *Source) {
	core.Province = make(map[string]*Province)
	for n, p := range source.Data.Diseaseh5Shelf.AreaTree[0].Children {
		if core.Province[p.Name] == nil {
			core.Province[p.Name] = &Province{
				UpdateTimeVirus:   p.Total.Mtime,
				AccComfirm:        p.Total.Confirm,
				NowComfirm:        p.Total.NowConfirm,
				AddComfirm:        p.Today.Confirm,
				AddComfirmLocal:   p.Today.LocalConfirmAdd,
				AddComfirmAboard:  p.Today.AbroadConfirmAdd,
				AccCure:           p.Total.Heal,
				AccDead:           p.Total.Dead,
				AddDead:           p.Today.DeadAdd,
				NowLocalAsymptoma: p.Total.Wzz,
				AddLocalAsymptoma: p.Today.WzzAdd,
				HighRiskAreaNum:   p.Total.HighRiskAreaNum,
				MediumRiskAreaNum: p.Total.MediumRiskAreaNum,
			}
		}

		// 地区解析
		for _, a := range source.Data.Diseaseh5Shelf.AreaTree[0].Children[n].Children {
			if core.Province[p.Name].Area == nil {
				core.Province[p.Name].Area = make(map[string]*Area)
			}
			core.Province[p.Name].Area[a.Name] = &Area{
				UpdateTimeVirus:   a.Total.Mtime,
				AccComfirm:        a.Total.Confirm,
				AddLocalComfirm:   a.Today.Confirm,
				AddLocalAsymptoma: a.Today.Confirm,
				HighRiskAreaNum:   a.Total.HighRiskAreaNum,
				MediumRiskAreaNum: a.Total.MediumRiskAreaNum,
			}
		}
	}
}

// 具体地区数据更新
func (core *Core) ParserArea(source *Source) {
	for _, a := range source.Data.LocalCityNCOVDataList {
		if core.Province[a.Province].Area == nil {
			core.Province[a.Province].Area = make(map[string]*Area)
		}
		area := core.Province[a.Province].Area[a.City]
		if area == nil {
			area = &Area{}
		}

		area.UpdateTimeVirus = a.Mtime
		dig, err := strconv.Atoi(a.LocalWzzAdd)
		if err != nil {
			dig = a.LocalConfirmAdd
		}
		area.AddLocalAsymptoma = dig
		area.AddLocalComfirm = a.LocalConfirmAdd
		area.HighRiskAreaNum = a.HighRiskAreaNum
		area.MediumRiskAreaNum = a.MediumRiskAreaNum
		// 回写
		core.Province[a.Province].Area[a.City] = area
	}
}

// 将输入的序列List转换成InlineKeyboard的序列
//
// label 作为一号位标识符
func (acore *Core) MakeInlineKeyboard(list [][]string, label string, isArea bool, province string) []tgbotapi.InlineKeyboardMarkup {
	core := []tgbotapi.InlineKeyboardButton{}
	ccore := [][]tgbotapi.InlineKeyboardButton{}
	cccore := []tgbotapi.InlineKeyboardMarkup{}

	page := 0
	rows := 5
	columns := 4
	row := 0
	col := 0
	Next := "» Next"
	Back := "« Back"
	Lock := false

	for {
		if page <= 0 {
			if len(list) > rows*columns-((row/1)*columns+col) {
				if row+1 == rows && col+1 == columns {
					if isArea {
						core = append(core, tgbotapi.NewInlineKeyboardButtonData(Next, fmt.Sprintf("%v-%v--%v-", label, page+1, province)))
					} else {
						core = append(core, tgbotapi.NewInlineKeyboardButtonData(Next, fmt.Sprintf("%v-%v-", label, page+1)))
					}
					col++
				} else {
					core = append(core, tgbotapi.NewInlineKeyboardButtonData(list[0][0], fmt.Sprintf("%v--%v", label, list[0][1])))
					if page == 0 && row == 0 && col == 0 && !isArea {
						if list[0][0] != "总览" {
							acore.Province[list[0][0]].PageNum = page
						}
					}
					col++
					list = list[1:]
				}
			} else {
				core = append(core, tgbotapi.NewInlineKeyboardButtonData(list[0][0], fmt.Sprintf("%v--%v", label, list[0][1])))
				if !isArea {
					if list[0][0] != "总览" && list[0][0] != Back && list[0][0] != Next && list[0][0] != "地区待确认" {
						acore.Province[list[0][0]].PageNum = page
					}
				}
				col++
				list = list[1:]
			}
		} else {
			if len(list) > rows*columns-(row/1*columns+col) {
				if row+1 == rows && col == 0 {
					if isArea {
						core = append(core, tgbotapi.NewInlineKeyboardButtonData(Back, fmt.Sprintf("%v-%v--%v-", label, page-1, province)))
					} else {
						core = append(core, tgbotapi.NewInlineKeyboardButtonData(Back, fmt.Sprintf("%v-%v-", label, page-1)))
					}
					col++
				} else if row+1 == rows && col+1 == columns {
					if isArea {
						core = append(core, tgbotapi.NewInlineKeyboardButtonData(Next, fmt.Sprintf("%v-%v--%v-", label, page+1, province)))
					} else {
						core = append(core, tgbotapi.NewInlineKeyboardButtonData(Next, fmt.Sprintf("%v-%v-", label, page+1)))
					}
					col++
				} else {
					core = append(core, tgbotapi.NewInlineKeyboardButtonData(list[0][0], fmt.Sprintf("%v--%v", label, list[0][1])))
					if !isArea {
						if list[0][0] != "总览" {
							acore.Province[list[0][0]].PageNum = page
						}
					}
					col++
					list = list[1:]
				}
			} else {
				if (len(list)+(row/1*columns+col))/columns == row && col == 0 {
					if isArea {
						core = append(core, tgbotapi.NewInlineKeyboardButtonData(Back, fmt.Sprintf("%v-%v--%v-", label, page-1, province)))
					} else {
						core = append(core, tgbotapi.NewInlineKeyboardButtonData(Back, fmt.Sprintf("%v-%v-", label, page-1)))
					}
					col++
					for _, i := range list {
						core = append(core, tgbotapi.NewInlineKeyboardButtonData(i[0], fmt.Sprintf("%v--%v", label, i[1])))
						if !isArea {
							if i[0] != "总览" {
								acore.Province[i[0]].PageNum = page
							}
						}
					}
					list = list[len(list):]
					Lock = true
				} else {
					core = append(core, tgbotapi.NewInlineKeyboardButtonData(list[0][0], fmt.Sprintf("%v--%v", label, list[0][1])))
					if !isArea {
						if list[0][0] != "总览" {
							acore.Province[list[0][0]].PageNum = page
						}
					}
					col++
					list = list[1:]
				}
			}
		}

		if len(core) == columns || len(list) == 0 {
			ccore = append(ccore, core)
			row++
			col = 0
			core = []tgbotapi.InlineKeyboardButton{}
			if !Lock && len(list) == 0 {
				if isArea {
					core = append(core, tgbotapi.NewInlineKeyboardButtonData(Back, fmt.Sprintf("%v-%v--%v-", label, page-1, province)))
				} else {
					core = append(core, tgbotapi.NewInlineKeyboardButtonData(Back, fmt.Sprintf("%v-%v-", label, page-1)))
				}
				ccore = append(ccore, core)
				core = []tgbotapi.InlineKeyboardButton{}
			}
		}

		if row/1*columns+col == rows*columns || len(list) == 0 {
			cccore = append(cccore, tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: ccore,
			})
			ccore = [][]tgbotapi.InlineKeyboardButton{}
			row, col = 0, 0
			page++
			if len(list) == 0 {
				break
			}
		}
	}
	return cccore
}

// 一级keyboard生成
func (core *Core) ProvinceKeyboard() {
	tempList := [][]string{
		{"总览", "pre"},
	}
	for k := range core.Province {
		tempList = append(tempList, []string{k, k})
	}
	core.ProvinceInlineKeyborad = core.MakeInlineKeyboard(tempList, "virus", false, "")
}

// 二级keyboard
func (core *Core) AreaKeyboard() {
	core.AreaInlineKeyboard = make(map[string][]tgbotapi.InlineKeyboardMarkup)
	for k := range core.Province {
		tempList := [][]string{
			{"« Back", fmt.Sprintf("back--%v", core.Province[k].PageNum)},
			{"总览", fmt.Sprintf("pre-%v-", k)},
		}
		for a := range core.Province[k].Area {
			if a != "境外输入" && a != "地区待确认" {
				tempList = append(tempList, []string{a, fmt.Sprintf("%v-%v-", a, k)})
			}
		}
		core.AreaInlineKeyboard[k] = core.MakeInlineKeyboard(tempList, "virus", true, k)
	}
}

// 具体地区的百度查询
func (paimon *Paimon) BaiduAreaQuery(province string, city string, core *Core) *Core {
	res, err := requests.Bronya("GET", fmt.Sprintf("https://opendata.baidu.com/data/inner?resource_id=5653&query=%v%v新型肺炎最新动态&alr=1", province, city), nil, nil, nil, false)
	if err != nil {
		paimon.Log.ERROR(err)
		return core
	}
	js := map[string]interface{}{}
	err = json.Unmarshal(res.Body, &js);
	if err != nil {
		paimon.Log.ERROR(err)
		return core
	}

	list := js["Result"].([]interface{})[0].(map[string]interface{})["DisplayData"].(map[string]interface{})["resultData"].(map[string]interface{})["tplData"].(map[string]interface{})["data_list"].([]interface{})

	area := core.Province[province].Area[city]
	for _, tli := range list {
		li := tli.(map[string]interface{})
		switch li["total_desc"].(string) {
		case "新增本土":
			{
				num, err := strconv.Atoi(li["total_num"].(string))
				if err != nil {
					paimon.Log.ERROR(err)
					num = area.AddLocalComfirm
				}
				area.AddLocalComfirm = num
			}
		case "新增本土无症状":
			{
				num, err := strconv.Atoi(li["total_num"].(string))
				if err != nil {
					paimon.Log.ERROR(err)
					num = area.AddLocalAsymptoma
				}
				area.AddLocalAsymptoma = num
			}
		case "现有确诊":
			{
				num, err := strconv.Atoi(li["total_num"].(string))
				if err != nil {
					paimon.Log.ERROR(err)
					num = area.NowComfirm
				}
				area.NowComfirm = num
			}
		case "累计确诊":
			{
				num, err := strconv.Atoi(li["total_num"].(string))
				if err != nil {
					paimon.Log.ERROR(err)
					num = area.AccComfirm
				}
				area.AccComfirm = num
			}
		case "累计治愈":
			{
				num, err := strconv.Atoi(li["total_num"].(string))
				if err != nil {
					paimon.Log.ERROR(err)
					num = area.AccCure
				}
				area.AccCure = num
			}
		case "累计死亡":
			{
				num, err := strconv.Atoi(li["total_num"].(string))
				if err != nil {
					paimon.Log.ERROR(err)
					num = area.AccDeadths
				}
				area.AccDeadths = num
			}
		}
	}
	core.Province[province].Area[city] = area
	return core
}

func (core *Core) GetPreChina() string {
	cn := core.China
	ctx := fmt.Sprintf(
		`国内总览
%v
累计确诊病例：%v
 └ 现有确诊病例：%v
 └ 新增确诊：%v
 └ 现有本土确诊：%v
 └ 新增本土确诊：%v
重症病例：%v
累计治愈：%v
累计死亡：%v
 └ 新增死亡：%v
现有本土无症状：%v
 └ 新增本土无症状：%v

%v
高风险地区：%v
中风险地区：%v`,
		cn.UpdateTimeVirus,
		cn.AccComfirm,
		cn.NowComfirm,
		cn.AddComfirm,
		cn.NowComfirmLocal,
		cn.AddComfirmLocal,
		cn.SevereCase,
		cn.AccCure,
		cn.AccDead,
		cn.AddDead,
		cn.NowLocalAsymptoma,
		cn.AddLocalAsymptoma,
		cn.UpdateTimeRisk,
		cn.HighRiskAreaNum,
		cn.MediumRiskAreaNum,
	)
	return ctx
}

func (core *Core) GetPreProvince(province string) string {
	pv := core.Province[province]
	ctx := fmt.Sprintf(
		`%v总览
%v
累计确诊病例：%v
 └ 现有确诊病例：%v
  └ 新增确诊：%v
   └ 新增本土确诊：%v
   └ 新增境外输入确诊：%v
累计治愈：%v
累计死亡：%v
 └ 新增死亡：%v
现有本土无症状：%v
 └ 新增本土无症状：%v

高风险地区：%v
中风险地区：%v`,
		province,
		pv.UpdateTimeVirus,
		pv.AccComfirm,
		pv.NowComfirm,
		pv.AddComfirm,
		pv.AddComfirmLocal,
		pv.AddComfirmAboard,
		pv.AccCure,
		pv.AccDead,
		pv.AddDead,
		pv.NowLocalAsymptoma,
		pv.AddLocalAsymptoma,
		pv.HighRiskAreaNum,
		pv.MediumRiskAreaNum,
	)
	return ctx
}

func (paimon *Paimon) GetArea(province string, area string, core *Core) string {
	core = paimon.BaiduAreaQuery(province, area, core)
	ar := core.Province[province].Area[area]
	ctx := fmt.Sprintf(
		`%v-%v
%v
累计确诊：%v
现有确诊：%v
 └ 新增本土确诊：%v
新增本土无症状：%v
累计治愈：%v
累计死亡：%v

高风险地区：%v
中风险地区：%v`,
		province, area,
		ar.UpdateTimeVirus,
		ar.AccComfirm,
		ar.NowComfirm,
		// 这个好像有问题
		// ar.AddComfirm,
		ar.AddLocalComfirm,
		ar.AddLocalAsymptoma,
		ar.AccCure,
		ar.AccDeadths,

		
		ar.HighRiskAreaNum,
		ar.MediumRiskAreaNum,
	)
	return ctx
}

// 入口
func MainHandle() (*Core, error) {
	log.SetFlags(log.Lshortfile)
	source, err := GetData()
	if err != nil {
		return nil, err
	}
	core := Core{}

	core.ChinaWin(source)
	// log.Println(core)
	core.ParseProvince(source)
	// log.Println(core)
	core.ParserArea(source)

	core.ProvinceKeyboard()
	// log.Println(core)
	core.AreaKeyboard()
	// log.Println(core)

	return &core, nil
}
