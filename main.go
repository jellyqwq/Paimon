// Paimon_poi_test_bot
// https://core.telegram.org/bots/api#using-a-local-bot-api-server
package main

import (
	"fmt"
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

type QueueInfo struct {
	TimeStamp int64
	MessageID int
	Core      *coronavirus.Core
}

var (
	log = &olog.Olog{
		Level: olog.LEVEL_ERROR,
	}
	CoronavirusQueue = make(map[int64]*QueueInfo)

	compileInlineInput = regexp.MustCompile(`^(?P<inlineType>.*?) +(?P<text>.*)`)
	compileElysia      = regexp.MustCompile(`^(派蒙|Paimon|飞行矮堇瓜|应急食品|白飞飞|神之嘴){1}`)

	// Inline keyboard
	HotwordKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("B站", "HotWordBilibili"),
			tgbotapi.NewInlineKeyboardButtonData("微博", "HotWordWeibo"),
		),
	)
	MusicSendKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:                         "y2mate.tools",
				SwitchInlineQueryCurrentChat: StringToPString("m1 "),
			},
			tgbotapi.InlineKeyboardButton{
				Text:                         "y2mate.com",
				SwitchInlineQueryCurrentChat: StringToPString("m2 "),
			},
		),
	)
)

func StringToPString(s string) *string {
	return &s
}

func deleteMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int, delay int64) {
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)
	time.Sleep(time.Duration(delay) * time.Second)
	bot.Send(msg)
}

func InitMessage(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	msg.DisableNotification = true
	return msg
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
						Core, err := coronavirus.MainHandle()
						if err != nil {
							log.ERROR(err)
							continue
						}
						if Core == nil {
							log.ERROR("Core is nil")
							continue
						}

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "疫情查询(￣_￣|||)")

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

						log.INFO(Core.ProvinceInlineKeyborad[0])
						msg.ReplyMarkup = Core.ProvinceInlineKeyborad[0]

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
				case "hoyocos":
					{
						list, err := webapi.HoyoBBS()
						if err != nil {
							log.ERROR(err)
							continue
						}
						var ImageList []interface{}
						for _, value := range list {
							Image := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(value))
							ImageList = append(ImageList, Image)
						}
						msg := tgbotapi.NewMediaGroup(update.Message.Chat.ID, ImageList)
						msg.DisableNotification = true
						go deleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID, config.DeleteMessageAfterSeconds)
						_, err = bot.Send(msg)
						if err != nil {
							log.ERROR(err)
							continue
						}
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
					core := CoronavirusQueue[cid].Core

					switch len(options) {
					case 3:
						{
							// virus page provice
							if options[2] == "pre" {
								// 总览操作
								msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, core.GetPreChina())
								msg = InitMessage(msg)
								if msg.Text != "" {
									if _, err := bot.Send(msg); err != nil {
										log.ERROR(err)
										return
									}
								}
							} else if options[1] != "" {
								// 省份页翻页操作, 不为空是翻页按钮
								index, err := strconv.Atoi(options[1])
								if err != nil {
									log.ERROR(err)
									continue
								}
								msg := tgbotapi.NewEditMessageReplyMarkup(cid, mid, core.ProvinceInlineKeyborad[index])
								if _, err := bot.Send(msg); err != nil {
									log.ERROR(err)
									continue
								}
							} else {
								// 其余情况是进入二级目录
								province := options[2]
								msg := tgbotapi.NewEditMessageReplyMarkup(cid, mid, core.AreaInlineKeyboard[province][0])
								if _, err := bot.Send(msg); err != nil {
									log.ERROR(err)
									continue
								}
							}

						}
					case 5:
						{
							// virus page (area|pre|back) Province ProvincePageNum
							if options[2] == "pre" {
								// 总览操作
								msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, core.GetPreProvince(options[3]))
								msg = InitMessage(msg)
								if msg.Text != "" {
									if _, err := bot.Send(msg); err != nil {
										log.ERROR(err)
										return
									}
								}
							} else if options[2] == "back" {
								// 返回上一级键盘
								index, err := strconv.Atoi(options[4])
								if err != nil {
									log.ERROR(err)
									continue
								}
								msg := tgbotapi.NewEditMessageReplyMarkup(cid, mid, core.ProvinceInlineKeyborad[index])
								if _, err := bot.Send(msg); err != nil {
									log.ERROR(err)
									continue
								}
							} else if options[1] != "" {
								// 地区页翻页操作, 不为空是翻页按钮
								index, err := strconv.Atoi(options[1])
								if err != nil {
									log.ERROR(err)
									continue
								}
								msg := tgbotapi.NewEditMessageReplyMarkup(cid, mid, core.AreaInlineKeyboard[options[3]][index])
								if _, err := bot.Send(msg); err != nil {
									log.ERROR(err)
									continue
								}
							} else {
								// 其余情况发送消息
								province := options[3]
								area := options[2]
								msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, core.GetArea(province, area))
								msg = InitMessage(msg)
								if _, err := bot.Send(msg); err != nil {
									log.ERROR(err)
									continue
								}
							}
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
