// Paimon_poi_test_bot
// https://core.telegram.org/bots/api#using-a-local-bot-api-server
package main

import (
	"fmt"
	// "log"
	"net/http"
	"strconv"
	"time"

	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/jellyqwq/Paimon/config"
	"github.com/jellyqwq/Paimon/coronavirus"
	"github.com/jellyqwq/Paimon/cqtotg"
	"github.com/jellyqwq/Paimon/news"
	"github.com/jellyqwq/Paimon/olog"
	"github.com/jellyqwq/Paimon/tools"
	"github.com/jellyqwq/Paimon/webapi"
)

var (
	log = &olog.Olog{
		Level: olog.LEVEL_ERROR,
	}
	CoronavirusQueue   = make(map[int64]*QueueInfo)
	compileInlineInput = regexp.MustCompile(`^(?P<inlineType>.*?) +(?P<text>.*)`)
	compileElysia = regexp.MustCompile(`^(派蒙|Paimon|飞行矮堇瓜|应急食品|白飞飞|神之嘴){1}`)
	// bot *tgbotapi.BotAPI

	// Inline keyboard
	HotwordKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("B站", "HotWordBilibili"),
			tgbotapi.NewInlineKeyboardButtonData("微博", "HotWordWeibo"),
		),
	)
	stringM1          = "m1 "
	stringM2          = "m2 "
	MusicSendKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:                         "y2mate.tools",
				SwitchInlineQueryCurrentChat: &stringM1,
			},
			tgbotapi.InlineKeyboardButton{
				Text:                         "y2mate.com",
				SwitchInlineQueryCurrentChat: &stringM2,
			},
		),
	)
)

type QueueInfo struct {
	TimeStamp      int64
	MessageID      int
	Core           *coronavirus.KernelVirus
	InlineKeyboard []tgbotapi.InlineKeyboardMarkup
}

func deleteMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int, delay int64) {
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)
	time.Sleep(time.Duration(delay) * time.Second)
	bot.Send(msg)
}

func mainHandler() {
	log.Update()
	config, err := config.ReadYaml()
	if err != nil {
		log.FATAL(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.FATAL(err)
	}

	bot.Debug = true

	// log.Printf("Authorized on account %s", bot.Self.UserName)

	// webhook, _ := tgbotapi.NewWebhookWithCert(config.Webhook.URL + bot.Token, tgbotapi.FilePath(config.Webhook.CertificatePemPath))
	webhook, _ := tgbotapi.NewWebhook(config.TelegramWebHook.Url + bot.Token)
	webhook.IPAddress = config.TelegramWebHook.IPAddress
	webhook.AllowedUpdates = config.TelegramWebHook.AllowedUpdates
	webhook.MaxConnections = config.TelegramWebHook.MaxConnections

	if _, err = bot.Request(webhook); err != nil {
		log.FATAL(err)
	}
	
	// cqhttp http-reverse
	botSet := &cqtotg.PostParams{Bot: bot, Conf: config}
	http.HandleFunc("/cq/", botSet.Post)

	// QQ video format server
	http.HandleFunc("/retype/", cqtotg.FileParse)

	// Y2mate by y2mate.tools
	params := &webapi.Params{Bot: bot, Conf: config}
	http.HandleFunc("/y2mate/tools/", params.Y2mateByTools)

	// Y2mate by y2mate.com
	http.HandleFunc("/y2mate/com/", params.Y2mateByCom)

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe(config.WebhookIP+":"+strconv.FormatUint(config.WebhookPort, 10), nil)

	for update := range updates {
		if time.Now().Local().Day() != log.Day {
			log.Update()
		}

		if update.Message != nil {

			text := update.Message.Text
			

			// inline keyboard with command
			if update.Message.IsCommand() {
				log.DEBUG(update.Message.Command())

				switch update.Message.Command() {
				case "hot_word":
					{
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "热搜🔥 | 大瓜🍉")
						msg.ReplyMarkup = HotwordKeyboard
						msg.DisableNotification = true

						rep, err := bot.Send(msg)
						if err != nil {
							log.ERROR(err)
							continue
						}

						go deleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID, config.DeleteMessageAfterSeconds)
						go deleteMessage(bot, rep.Chat.ID, rep.MessageID, config.DeleteMessageAfterSeconds)
					}
				case "finance":
					{
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🏦💰货币汇率💸")

						CurrencyList := config.Currency

						var ResultList []string
						for n, s := range CurrencyList {
							for m, e := range CurrencyList {
								if n != m {
									ResultList = append(ResultList, s+"-"+e)
								}
							}
						}

						var keyboard [][]tgbotapi.InlineKeyboardButton
						var row []tgbotapi.InlineKeyboardButton
						var c int = 1
						for _, li := range ResultList {
							row = append(row, tgbotapi.NewInlineKeyboardButtonData(li, fmt.Sprintf("currency-%v", li)))
							// 每四个块合并row到keyboard中并重置row
							if c%3 == 0 {
								keyboard = append(keyboard, row)
								row = nil
								c = 0
							}
							c += 1
						}

						msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
							InlineKeyboard: keyboard,
						}
						msg.DisableNotification = true

						rep, err := bot.Send(msg)
						if err != nil {
							log.ERROR(err)
							continue
						}

						go deleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID, config.DeleteMessageAfterSeconds)
						go deleteMessage(bot, rep.Chat.ID, rep.MessageID, config.DeleteMessageAfterSeconds)
					}
				case "help":
					{
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "[Paimon | 应急食品](https://github.com/jellyqwq/Paimon)\n1. *点歌* _@Paimon_poi_bot <m1|m2> music name_ (m1是[y2mate.tools](y2mate.tools) | m2是[y2mate.com](www.y2mate.com))\n2. *信息查看* _派蒙INFO_ (单独发或Reply)\n3. *翻译句子* _派蒙翻译_ (配上句子发或Reply)\n4. *Command*")
						msg.ParseMode = "Markdown"
						msg.DisableWebPagePreview = true
						msg.DisableNotification = true

						if _, err := bot.Send(msg); err != nil {
							log.ERROR(err)
							continue
						}

						go deleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID, config.DeleteMessageAfterSeconds)
					}
				case "music":
					{
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "油管白嫖🎼")
						msg.ReplyMarkup = MusicSendKeyboard
						msg.DisableNotification = true

						rep, err := bot.Send(msg)
						if err != nil {
							log.ERROR(err)
							continue
						}

						go deleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID, config.DeleteMessageAfterSeconds)
						go deleteMessage(bot, rep.Chat.ID, rep.MessageID, config.DeleteMessageAfterSeconds)
					}
				case "coronavirus":
					{
						Core, err := coronavirus.Entry()
						if err != nil {
							log.ERROR(err)
							continue
						}
						if Core == nil {
							log.ERROR("Core is nil")
							continue
						}
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%v\n%v", Core.Time, Core.Title))

						Province := []string{}
						for k := range Core.ProvinceDetailed {
							Province = append(Province, k)
						}

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
								if len(Province) > rows*columns-((row/1)*columns+col) {
									if row+1 == rows && col+1 == columns {
										core = append(core, tgbotapi.NewInlineKeyboardButtonData(Next, fmt.Sprintf("virus-%v-", page+1)))
										col++
									} else {
										core = append(core, tgbotapi.NewInlineKeyboardButtonData(Province[0], fmt.Sprintf("virus--%v", Province[0])))
										col++
										Province = Province[1:]
									}
								} else {
									core = append(core, tgbotapi.NewInlineKeyboardButtonData(Province[0], fmt.Sprintf("virus--%v", Province[0])))
									col++
									Province = Province[1:]
								}
							} else {
								if len(Province) > rows*columns-(row/1*columns+col) {
									if row+1 == rows && col == 0 {
										core = append(core, tgbotapi.NewInlineKeyboardButtonData(Back, fmt.Sprintf("virus-%v-", page-1)))
										col++
									} else if row+1 == rows && col+1 == columns {
										core = append(core, tgbotapi.NewInlineKeyboardButtonData(Next, fmt.Sprintf("virus-%v-", page+1)))
										col++
									} else {
										core = append(core, tgbotapi.NewInlineKeyboardButtonData(Province[0], fmt.Sprintf("virus--%v", Province[0])))
										col++
										Province = Province[1:]
									}
								} else {
									if (len(Province)+(row/1*columns+col))/columns == row && col == 0 {
										core = append(core, tgbotapi.NewInlineKeyboardButtonData(Back, fmt.Sprintf("virus-%v-", page-1)))
										col++
										for _, i := range Province {
											core = append(core, tgbotapi.NewInlineKeyboardButtonData(i, fmt.Sprintf("virus--%v", i)))
										}
										Province = Province[len(Province):]
										Lock = true
									} else {
										core = append(core, tgbotapi.NewInlineKeyboardButtonData(Province[0], fmt.Sprintf("virus--%v", Province[0])))
										col++
										Province = Province[1:]
									}
								}
							}

							if len(core) == columns || len(Province) == 0 {
								ccore = append(ccore, core)
								row++
								col = 0
								core = []tgbotapi.InlineKeyboardButton{}
								if !Lock && len(Province) == 0 {
									core = append(core, tgbotapi.NewInlineKeyboardButtonData(Back, fmt.Sprintf("virus-%v-", page-1)))
									ccore = append(ccore, core)
									core = []tgbotapi.InlineKeyboardButton{}
								}
							}

							if row/1*columns+col == rows*columns || len(Province) == 0 {
								cccore = append(cccore, tgbotapi.InlineKeyboardMarkup{
									InlineKeyboard: ccore,
								})
								ccore = [][]tgbotapi.InlineKeyboardButton{}
								row, col = 0, 0
								page++
								if len(Province) == 0 {
									break
								}
							}
						}

						chatID := update.Message.Chat.ID
						if CoronavirusQueue[chatID] == nil {
							CoronavirusQueue[chatID] = &QueueInfo{}
						}
						// 记录
						TempStruct := CoronavirusQueue[chatID]

						// 删除上一个keyboard
						if TempStruct.MessageID != 0 {
							go deleteMessage(bot, chatID, TempStruct.MessageID, 0)
						}
						// 更新消息id
						TempStruct.MessageID = update.Message.MessageID
						// 更新时间戳
						TempStruct.TimeStamp = time.Now().Unix()
						// 更新核心
						TempStruct.Core = Core
						// 更新keyboard list
						TempStruct.InlineKeyboard = cccore

						msg.ReplyMarkup = cccore[0]
						msg.DisableNotification = true
						res, err := bot.Send(msg)
						if err != nil {
							log.ERROR(err)
							continue
						}

						TempStruct.MessageID = res.MessageID
						CoronavirusQueue[chatID] = TempStruct

						go deleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID, config.DeleteMessageAfterSeconds)
					}
				}

			} else if compileElysia.Match([]byte(text)) {
				text = string(compileElysia.ReplaceAll([]byte(text), []byte("")))

				var msg tgbotapi.MessageConfig
				if strings.Contains(text, "INFO") {
					var ctx string
					if update.Message.ReplyToMessage != nil {
						ctx = fmt.Sprintf("ReplyUserInfo\nUserID:`%v`\nChatID:`%v`\nFirstName:`%v`\nLastName:`%v`\nUserName:`%v`", update.Message.ReplyToMessage.From.ID, update.Message.ReplyToMessage.Chat.ID, update.Message.ReplyToMessage.From.FirstName, update.Message.ReplyToMessage.From.LastName, update.Message.ReplyToMessage.From.UserName)
					} else {
						ctx = fmt.Sprintf("UserInfo\nUserID:`%v`\nChatID:`%v`\nFirstName:`%v`\nLastName:`%v`\nUserName:`%v`", update.Message.From.ID, update.Message.Chat.ID, update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName)
					}

					msg = tgbotapi.NewMessage(update.Message.Chat.ID, ctx)
					msg.ParseMode = "Markdown"
					msg.DisableNotification = true
					msg.ReplyToMessageID = update.Message.MessageID

				} else if strings.Contains(text, "翻译") {
					compileTranslate := regexp.MustCompile(`翻译`)
					text := compileTranslate.ReplaceAllString(text, "")
					if update.Message.ReplyToMessage != nil {
						text = update.Message.ReplyToMessage.Text
					}
					if len(text) > 0 {
						result, err := webapi.RranslateByYouDao(text)
						if err != nil {
							log.ERROR(err)
							continue
						} else if len(result) > 0 {
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, result)
							msg.ReplyToMessageID = update.Message.MessageID
							msg.DisableNotification = true
						} else {
							continue
						}
					}

				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
				}
				msg.DisableNotification = true
				if msg.Text != "" {
					if _, err := bot.Send(msg); err != nil {
						log.ERROR(err)
						continue
					}
				}
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "loading...")
			if _, err := bot.Request(callback); err != nil {
				log.ERROR(err)
				continue
			}
			
			CallbackQueryData := update.CallbackQuery.Data
			if CallbackQueryData == "HotWordBilibili" {

				ctx, err := news.BiliHotWords()
				if err != nil {
					log.ERROR(err)
					continue
				}

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx)
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				msg.DisableNotification = true
				if msg.Text != "" {
					if _, err := bot.Send(msg); err != nil {
						log.ERROR(err)
						continue
					}
				}
			} else if CallbackQueryData == "HotWordWeibo" {
				ctx, err := news.WeiboHotWords()
				if err != nil {
					log.ERROR(err)
					continue
				}

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx)
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				msg.DisableNotification = true
				if msg.Text != "" {
					if _, err := bot.Send(msg); err != nil {
						log.ERROR(err)
						continue
					}
				}
			} else if len(CallbackQueryData) > 7 {
				if CallbackQueryData[:5] == "virus" {
					mid := update.CallbackQuery.Message.MessageID
					cid := update.CallbackQuery.Message.Chat.ID
					options := strings.Split(CallbackQueryData, "-")

					log.DEBUG(options)

					if options[2] != "" {
						SubCore := CoronavirusQueue[cid].Core.ProvinceDetailed[options[2]].New.Diagnose
						ctx := fmt.Sprintf("%v\n%v\n省份: %v\n新增境外输入: %v\n└无症状转确诊: %v\n新增本土: %v\n└无症状转确诊: %v", CoronavirusQueue[cid].Core.Time, CoronavirusQueue[cid].Core.Title, options[2], SubCore.Abroad, SubCore.AbroadFromAsymptoma, SubCore.Mainland, SubCore.MainlandFromAsymptoma)
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx)
						msg.ParseMode = "Markdown"
						msg.DisableWebPagePreview = true
						msg.DisableNotification = true
						if msg.Text != "" {
							if _, err := bot.Send(msg); err != nil {
								log.ERROR(err)
								return
							}
						}
					} else if options[1] != "" {
						index, err := strconv.Atoi(options[1])
						if err != nil {
							log.ERROR(err)
							continue
						}
						msg := tgbotapi.NewEditMessageReplyMarkup(cid, mid, CoronavirusQueue[cid].InlineKeyboard[index])
						if _, err := bot.Send(msg); err != nil {
							log.ERROR(err)
							continue
						}
					}

				} else if len(CallbackQueryData) > 10 {
					if CallbackQueryData[:8] == "currency" {
						tempList := strings.Split(CallbackQueryData, "-")
						currency := tempList[1] + "-" + tempList[2]
						ctx, err := webapi.Finance(currency)
						if err != nil {
							log.ERROR(err)
							continue
						}
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx)
						msg.ParseMode = "Markdown"
						msg.DisableWebPagePreview = true
						msg.DisableNotification = true
						if msg.Text != "" {
							if _, err := bot.Send(msg); err != nil {
								log.ERROR(err)
								continue
							}
						}
					}
				}
			}
		} else if update.InlineQuery != nil {
			text := update.InlineQuery.Query
			params := &webapi.Params{Bot: bot, Conf: config}

			if len(text) > 3 {

				paramsMap := tools.GetParamsOneDimension(compileInlineInput, text)
				inlineType := paramsMap["inlineType"]
				text := paramsMap["text"]
				if len(text) == 0 {
					continue
				}

				result := []interface{}{}
				// m1是y2mate.tools m2是www.y2mate.com
				switch inlineType {
				case "m1":
					{
						result, err = params.YoutubeSearch(text, inlineType)
						if err != nil {
							log.ERROR(err)
							continue
						}
					}
				case "m2":
					{
						result, err = params.YoutubeSearch(text, inlineType)
						if err != nil {
							log.ERROR(err)
							continue
						}
					}
				}

				if len(result) == 0 {
					continue
				}

				inlineConf := tgbotapi.InlineConfig{
					InlineQueryID: update.InlineQuery.ID,
					IsPersonal:    false,
					CacheTime:     0,
					Results:       result,
				}
				if _, err := bot.Send(inlineConf); err != nil {
					log.ERROR(err)
					continue
				}
			}
		}
	}
}

func main() {
	mainHandler()
}
