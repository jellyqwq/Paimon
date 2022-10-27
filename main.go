// Paimon_poi_test_bot
// https://core.telegram.org/bots/api#using-a-local-bot-api-server
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"regexp"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/jellyqwq/Paimon/config"
	"github.com/jellyqwq/Paimon/cqtotg"
	"github.com/jellyqwq/Paimon/news"
	"github.com/jellyqwq/Paimon/tools"

	// "github.com/jellyqwq/Paimon/requests"
	"github.com/jellyqwq/Paimon/webapi"
)

// Inline keyboard
var HotwordKeyboard = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("B站", "HotWordBilibili"),
        tgbotapi.NewInlineKeyboardButtonData("微博", "HotWordWeibo"),
    ),
)

var FinanceKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("USD-CNY", "USD-CNY"),
		tgbotapi.NewInlineKeyboardButtonData("CNY-USD", "CNY-USD"),
	),
)

func logError(err error) {
	if err != nil {
	  	log.Println(err)
	}
}

func mainHandler() {
	// 全局作用的正则表达式编译
	compileInlineInput := regexp.MustCompile(`^(?P<inlineType>.*?) +(?P<text>.*)`)

	config, err := config.ReadYaml()
	logError(err)

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	logError(err)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// webhook, _ := tgbotapi.NewWebhookWithCert(config.Webhook.URL + bot.Token, tgbotapi.FilePath(config.Webhook.CertificatePemPath))
	webhook, _ := tgbotapi.NewWebhook(config.TelegramWebHook.Url + bot.Token)
	webhook.IPAddress = config.TelegramWebHook.IPAddress
	webhook.AllowedUpdates = config.TelegramWebHook.AllowedUpdates
	webhook.MaxConnections = config.TelegramWebHook.MaxConnections

	_, err = bot.Request(webhook)
	logError(err)

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
	go http.ListenAndServe(config.WebhookIP + ":" + strconv.FormatUint(config.WebhookPort, 10), nil)

	
	for update := range updates {

		if update.Message != nil {

			text := update.Message.Text
			regElysia := regexp.MustCompile(`^(派蒙|Paimon|飞行矮堇瓜|应急食品|白飞飞|神之嘴)?`)

			// inline keyboard with command
			if update.Message.IsCommand() {
				log.Println(update.Message.Command())
				switch update.Message.Command() {
					case "hot_word": {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "热搜🔥 | 大瓜🍉")
						msg.ReplyMarkup = HotwordKeyboard
						msg.DisableNotification = true
						if _, err = bot.Send(msg); err != nil {
							panic(err)
						}
						deletmsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
						bot.Send(deletmsg)
					}
					case "finance": {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🏦💰货币汇率💸")
						msg.ReplyMarkup = FinanceKeyboard
						msg.DisableNotification = true
						if _, err = bot.Send(msg); err != nil {
							panic(err)
						}
						deletmsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
						bot.Send(deletmsg)
					}
					case "help": {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "[Paimon | 应急食品](https://github.com/jellyqwq/Paimon)\n1. *点歌* _@Paimon_poi_bot <m1|m2> music name_ (m1是[y2mate.tools](y2mate.tools) | m2是[y2mate.com](www.y2mate.com))\n2. *信息查看* _派蒙INFO_ (单独发或Reply)\n3. *翻译句子* _派蒙翻译_ (配上句子发或Reply)\n4. *Command*")
						msg.ParseMode = "Markdown"
						msg.DisableWebPagePreview = true
						msg.DisableNotification = true
						if _, err = bot.Send(msg); err != nil {
							panic(err)
						}
						deletmsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
						bot.Send(deletmsg)
					}
				}
				
			} else if (regElysia.Match([]byte(text))) {
				text = string(regElysia.ReplaceAll([]byte(text), []byte("")))

				var msg tgbotapi.MessageConfig
				if strings.Contains(text, "INFO"){
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
							log.Println(err)
							continue
						} else if len(result) > 0 {
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, result)
							msg.ReplyToMessageID = update.Message.MessageID
							msg.DisableNotification = true
						} else {
							continue
						}
					}
					
				} else{
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
				}
				msg.DisableNotification = true
				if msg.Text != "" {
					_, err := bot.Send(msg)
					logError(err)
				}
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}
			if callback.Text == "HotWordBilibili" {
				ctx, err := news.BiliHotWords()
				logError(err)

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx)
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				msg.DisableNotification = true
				if msg.Text != "" {
					_, err := bot.Send(msg)
					logError(err)
				}
			} else if callback.Text == "HotWordWeibo" {
				ctx, err := news.WeiboHotWords()
				logError(err)

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx)
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				msg.DisableNotification = true
				if msg.Text != "" {
					_, err := bot.Send(msg)
					logError(err)
				}
			} else if callback.Text == "USD-CNY" || callback.Text == "CNY-USD" {
				ctx, err := webapi.Finance(callback.Text)
				logError(err)

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx)
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				msg.DisableNotification = true
				if msg.Text != "" {
					_, err := bot.Send(msg)
					logError(err)
				}
			}
		} else if update.InlineQuery != nil {
			text := update.InlineQuery.Query
			params := &webapi.Params{Bot: bot, Conf: config}
			log.Println(text, len(text))

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
					case "m1": {
						result, err = params.YoutubeSearch(text, inlineType)
						if err != nil {
							log.Println(err)
							continue
						}
					}
					case "m2": {
						result, err = params.YoutubeSearch(text, inlineType)
						if err != nil {
							log.Println(err)
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
				bot.Send(inlineConf)
			}
		}
	}
}

func main() {
	mainHandler()
}