package coronavirus

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Source struct {
	Ret  int    `json:"ret"`
	Info string `json:"info"`
	Data struct {
		Diseaseh5Shelf struct {
			ChinaAdd struct {
				Heal           int `json:"heal"`
				Dead           int `json:"dead"`
				NowConfirm     int `json:"nowConfirm"`
				Suspect        int `json:"suspect"`
				NoInfect       int `json:"noInfect"`
				LocalConfirmH5 int `json:"localConfirmH5"`
				Confirm        int `json:"confirm"`
				NowSevere      int `json:"nowSevere"`
				ImportedCase   int `json:"importedCase"`
				LocalConfirm   int `json:"localConfirm"`
				NoInfectH5     int `json:"noInfectH5"`
			} `json:"chinaAdd"`
			IsShowAdd     bool `json:"isShowAdd"`
			ShowAddSwitch struct {
				Confirm        bool `json:"confirm"`
				NowConfirm     bool `json:"nowConfirm"`
				NowSevere      bool `json:"nowSevere"`
				NoInfect       bool `json:"noInfect"`
				LocalConfirm   bool `json:"localConfirm"`
				All            bool `json:"all"`
				Suspect        bool `json:"suspect"`
				Dead           bool `json:"dead"`
				Heal           bool `json:"heal"`
				ImportedCase   bool `json:"importedCase"`
				Localinfeciton bool `json:"localinfeciton"`
			} `json:"showAddSwitch"`
			AreaTree []struct {
				Children []struct {
					Name   string `json:"name"`
					Adcode string `json:"adcode"`
					Date   string `json:"date"`
					Today  struct {
						ConfirmCuts      int    `json:"confirmCuts"`
						IsUpdated        bool   `json:"isUpdated"`
						Tip              string `json:"tip"`
						WzzAdd           int    `json:"wzz_add"`
						LocalConfirmAdd  int    `json:"local_confirm_add"`
						AbroadConfirmAdd int    `json:"abroad_confirm_add"`
						DeadAdd          int    `json:"dead_add"`
						Confirm          int    `json:"confirm"`
					} `json:"today"`
					Total struct {
						NowConfirm                     int    `json:"nowConfirm"`
						Dead                           int    `json:"dead"`
						ShowHeal                       bool   `json:"showHeal"`
						Wzz                            int    `json:"wzz"`
						ProvinceLocalConfirm           int    `json:"provinceLocalConfirm"`
						ContinueDayZeroConfirm         int    `json:"continueDayZeroConfirm"`
						Mtime                          string `json:"mtime"`
						ShowRate                       bool   `json:"showRate"`
						HighRiskAreaNum                int    `json:"highRiskAreaNum"`
						ContinueDayZeroLocalConfirmAdd int    `json:"continueDayZeroLocalConfirmAdd"`
						Adcode                         string `json:"adcode"`
						Confirm                        int    `json:"confirm"`
						Heal                           int    `json:"heal"`
						MediumRiskAreaNum              int    `json:"mediumRiskAreaNum"`
						ContinueDayZeroConfirmAdd      int    `json:"continueDayZeroConfirmAdd"`
					} `json:"total"`
					Children []struct {
						Name   string `json:"name"`
						Adcode string `json:"adcode"`
						Date   string `json:"date"`
						Today  struct {
							WzzAdd          string `json:"wzz_add"`
							LocalConfirmAdd int    `json:"local_confirm_add"`
							Confirm         int    `json:"confirm"`
							ConfirmCuts     int    `json:"confirmCuts"`
							IsUpdated       bool   `json:"isUpdated"`
						} `json:"today"`
						Total struct {
							HighRiskAreaNum                int    `json:"highRiskAreaNum"`
							Mtime                          string `json:"mtime"`
							ShowHeal                       bool   `json:"showHeal"`
							Heal                           int    `json:"heal"`
							MediumRiskAreaNum              int    `json:"mediumRiskAreaNum"`
							Dead                           int    `json:"dead"`
							ShowRate                       bool   `json:"showRate"`
							Wzz                            int    `json:"wzz"`
							ContinueDayZeroLocalConfirmAdd int    `json:"continueDayZeroLocalConfirmAdd"`
							Confirm                        int    `json:"confirm"`
							ProvinceLocalConfirm           int    `json:"provinceLocalConfirm"`
							ContinueDayZeroLocalConfirm    int    `json:"continueDayZeroLocalConfirm"`
							Adcode                         string `json:"adcode"`
							NowConfirm                     int    `json:"nowConfirm"`
						} `json:"total"`
					} `json:"children"`
				} `json:"children"`
				Name  string `json:"name"`
				Today struct {
					Confirm   int  `json:"confirm"`
					IsUpdated bool `json:"isUpdated"`
				} `json:"today"`
				Total struct {
					Mtime                          string `json:"mtime"`
					Adcode                         string `json:"adcode"`
					Dead                           int    `json:"dead"`
					ShowRate                       bool   `json:"showRate"`
					ProvinceLocalConfirm           int    `json:"provinceLocalConfirm"`
					ContinueDayZeroLocalConfirmAdd int    `json:"continueDayZeroLocalConfirmAdd"`
					ContinueDayZeroLocalConfirm    int    `json:"continueDayZeroLocalConfirm"`
					Confirm                        int    `json:"confirm"`
					Heal                           int    `json:"heal"`
					NowConfirm                     int    `json:"nowConfirm"`
					ShowHeal                       bool   `json:"showHeal"`
					Wzz                            int    `json:"wzz"`
					MediumRiskAreaNum              int    `json:"mediumRiskAreaNum"`
					HighRiskAreaNum                int    `json:"highRiskAreaNum"`
				} `json:"total"`
			} `json:"areaTree"`
			LastUpdateTime string `json:"lastUpdateTime"`
			ChinaTotal     struct {
				Heal               int    `json:"heal"`
				Dead               int    `json:"dead"`
				NowSevere          int    `json:"nowSevere"`
				ImportedCase       int    `json:"importedCase"`
				ConfirmAdd         int    `json:"confirmAdd"`
				DeadAdd            int    `json:"deadAdd"`
				Mtime              string `json:"mtime"`
				MediumRiskAreaNum  int    `json:"mediumRiskAreaNum"`
				Suspect            int    `json:"suspect"`
				ShowLocalConfirm   int    `json:"showLocalConfirm"`
				LocalAccConfirm    int    `json:"local_acc_confirm"`
				NowLocalWzz        int    `json:"nowLocalWzz"`
				NowConfirm         int    `json:"nowConfirm"`
				LocalConfirm       int    `json:"localConfirm"`
				Showlocalinfeciton int    `json:"showlocalinfeciton"`
				LocalConfirmH5     int    `json:"localConfirmH5"`
				LocalWzzAdd        int    `json:"localWzzAdd"`
				HighRiskAreaNum    int    `json:"highRiskAreaNum"`
				Confirm            int    `json:"confirm"`
				NoInfect           int    `json:"noInfect"`
				NoInfectH5         int    `json:"noInfectH5"`
				LocalConfirmAdd    int    `json:"localConfirmAdd"`
				MRiskTime          string `json:"mRiskTime"`
			} `json:"chinaTotal"`
		} `json:"diseaseh5Shelf"`
		LocalCityNCOVDataList []struct {
			Province          string `json:"province"`
			Adcode            string `json:"adcode"`
			Date              string `json:"date"`
			Mtime             string `json:"mtime"`
			HighRiskAreaNum   int    `json:"highRiskAreaNum"`
			City              string `json:"city"`
			IsUpdated         bool   `json:"isUpdated"`
			LocalConfirmAdd   int    `json:"local_confirm_add"`
			LocalWzzAdd       string `json:"local_wzz_add"`
			MediumRiskAreaNum int    `json:"mediumRiskAreaNum"`
			IsSpecialCity     bool   `json:"isSpecialCity"`
		} `json:"localCityNCOVDataList"`
	} `json:"data"`
}

type Area struct {
	// ????????????????????????
	UpdateTimeVirus string
	// ??????????????????
	NowComfirm int
	// ??????????????????
	AddComfirm int
	// ?????????????????????
	// NowLocalAsymptoma int
	// ?????????????????????
	AddLocalAsymptoma int
	// ???????????????
	HighRiskAreaNum int
	// ???????????????
	MediumRiskAreaNum int
}

type Province struct {
	// ?????????
	PageNum int
	// ????????????????????????
	UpdateTimeVirus string
	// ??????????????????
	AccComfirm int
	// ??????????????????
	NowComfirm int
	// ??????????????????
	AddComfirm int
	// ??????????????????
	AddComfirmLocal int
	// ????????????????????????
	AddComfirmAboard int
	// ????????????
	AccCure int
	// ????????????
	AccDead int
	// ????????????
	AddDead int
	// ?????????????????????
	NowLocalAsymptoma int
	// ?????????????????????
	AddLocalAsymptoma int
	// ???????????????
	HighRiskAreaNum int
	// ???????????????
	MediumRiskAreaNum int
	// ????????????
	Area map[string]*Area
}

// Core ???
//
// ???????????????????????????
type Core struct {
	// ????????????
	Province map[string]*Province
	China    struct {
		// ????????????????????????
		UpdateTimeVirus string
		// ??????????????????
		AccComfirm int
		// ??????????????????
		NowComfirm int
		// ??????????????????
		AddComfirm int
		// ??????????????????
		NowComfirmLocal int
		// ??????????????????
		AddComfirmLocal int
		// ????????????
		SevereCase int
		// ????????????
		AccCure int
		// ????????????
		AccDead int
		// ????????????
		AddDead int
		// ?????????????????????
		NowLocalAsymptoma int
		// ?????????????????????
		AddLocalAsymptoma int
		// ?????????????????????
		UpdateTimeRisk string
		// ???????????????
		HighRiskAreaNum int
		// ???????????????
		MediumRiskAreaNum int
	}
	ProvinceInlineKeyborad []tgbotapi.InlineKeyboardMarkup
	AreaInlineKeyboard     map[string][]tgbotapi.InlineKeyboardMarkup
}
