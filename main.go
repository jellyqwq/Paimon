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

	// "github.com/jellyqwq/Paimon/cqtotg"
	"github.com/jellyqwq/Paimon/news"
	"github.com/jellyqwq/Paimon/olog"
	"github.com/jellyqwq/Paimon/tools"
	"github.com/jellyqwq/Paimon/webapi"
	"github.com/chromedp/chromedp"
	"context"
)

var (
	log = &olog.Olog{
		Level: olog.LEVEL_DEBUG,
	}

	compileInlineInput = regexp.MustCompile(`^(?P<inlineType>.*?) +(?P<text>.*)`)
	compileElysia      = regexp.MustCompile(`^(派蒙|Paimon|飞行矮堇瓜|应急食品|白飞飞|神之嘴){1}`)
	// 米游社链接匹配
	compilMihousheArticle = regexp.MustCompile(`(www|m)\.miyoushe\.com.*?article/\d+`)

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

func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(".mhy-article-page__main img", chromedp.ByQueryAll),
		chromedp.WaitVisible("i.mhy-icon.iconfont.icon-liulanshu,.icon-liuyanshu,.icon-dianzan,.icon-shoucang", chromedp.ByQueryAll),
		chromedp.Sleep(4* time.Second),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
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
	// botSet := &cqtotg.PostParams{Bot: bot, Conf: config}
	// http.HandleFunc("/cq/", botSet.Post)

	// QQ video format server
	// http.HandleFunc("/retype/", cqtotg.FileParse)

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
				case "livecode":
					{
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, webapi.MihoyoLiveCode())
						msg.ParseMode = "Markdown"
						msg.DisableWebPagePreview = true
						msg.DisableNotification = true
						go deleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID, config.DeleteMessageAfterSeconds)

						if _, err := bot.Send(msg); err != nil {
							log.ERROR(err)
							continue
						}
					}
				}
				// } else if compileElysia.Match([]byte(text)) {
			} else if update.Message.Text != "" {
				text := update.Message.Text
				text = string(compileElysia.ReplaceAll([]byte(text), []byte("")))
				log.DEBUG(text)

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
				} else if compilMihousheArticle.MatchString(text) {
					url := "https://" + compilMihousheArticle.FindString(text)
					// 设置项，具体参照参考1
					options := []chromedp.ExecAllocatorOption{
						chromedp.Flag("headless", true),                      // debug使用
						chromedp.Flag("blink-settings", "imagesEnabled=true"), // 禁用图片加载
						chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
						// chromedp.UserAgent(`Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1 miHoYoBBS/2.41.2`),
					}
					options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

					// Chrome初始化代码如下：
					c, _ := chromedp.NewExecAllocator(context.Background(), options...)

					go func(update tgbotapi.Update) {
						dpctx, cancel := chromedp.NewContext(c)
						defer cancel()
						dpctx, cancel = context.WithTimeout(dpctx, 60*time.Second)
						defer cancel()

						var imageBuf []byte

						// 捕获某个元素的截图
						if err := chromedp.Run(dpctx,
							elementScreenshot(url, `div.mhy-article-page__main`, &imageBuf),
							// elementScreenshot(url, `div.mhy-card.main-article`, &imageBuf),
						); err != nil {
							log.ERROR(err)
							return
						}
						msg := tgbotapi.NewDocument(update.Message.Chat.ID, 
							tgbotapi.FileBytes{
								Name: url + ".png",
								Bytes: imageBuf,
							},
						)
						msg.DisableNotification = true
						msg.ReplyToMessageID = update.Message.MessageID
						if _, err := bot.Send(msg); err != nil {
							log.ERROR(err)
							log.INFO("尝试使用普通发送")
							msg := tgbotapi.NewPhoto(update.Message.Chat.ID, 
								tgbotapi.FileBytes{
									Name: url,
									Bytes: imageBuf,
								},
							)
							msg.DisableNotification = true
							msg.ReplyToMessageID = update.Message.MessageID
							if _, err := bot.Send(msg); err != nil {
								log.ERROR(err)
								return
							}
							return
						}
					}(update)
					continue
				} else if text == "原神" {
					ctx := "差不多得了😅屁大点事都要拐上原神，原神一没招你惹你，二没干伤天害理的事情，到底怎么你了让你一直无脑抹黑，米哈游每天费尽心思的文化输出弘扬中国文化，你这种喷子只会在网上敲键盘诋毁良心公司，中国游戏的未来就是被你这种人毁掉的😅\n叫我们原批的小心点\n老子在大街上亲手给打过两个\n我在公共座椅上无聊玩原神，有两个B就从我旁边过，看见我玩原神就悄悄说了一句:又是一个原批，我就直接上去一拳呼脸上，我根本不给他解释的机会，我也不问他为什么说我是原批，我就打，我就看他不爽，他惹我了，我就不给他解释的机会，直接照着脸和脑门就打直接给那B呼出鼻血，脸上青一块，紫一块的我没撕她嘴巴都算好了你们这还不算最狠的，我记得我以前小时候春节去老家里，有一颗核弹，我以为是鞭炮，和大地红一起点了，当时噼里啪啦得，然后突然一朵蘑菇云平地而起，当时我就只记得两眼一黑，昏过去了，整个村子没了，幸好我是体育生，身体素质不错，住了几天院就没事了，几个月下来腿脚才利落，现在已经没事了，但是那种钻心的疼还是让我一生难忘，😂😂😂  令人感叹今早一玩原神，我便昏死了过去，现在才刚刚缓过来。在昏死过去的短短数小时内，我的大脑仿佛被龙卷风无数次摧毁。\n在原神这一神作的面前，我就像一个一丝不挂的原始人突然来到了现代都市，二次元已如高楼大厦将我牢牢地吸引，开放世界就突然变成那喇叭轰鸣的汽车，不仅把我吓个措手不及，还让我瞬间将注意完全放在了这新的奇物上面，而还没等我稍微平复心情，纹化输出的出现就如同眼前遮天蔽日的宇宙战舰，将我的世界观无情地粉碎，使我彻底陷入了忘我的迷乱，狂泄不止。\n原神，那眼花缭乱的一切都让我感到震撼，但是我那贫瘠的大脑却根本无法理清其中任何的逻辑，巨量的信息和情感泄洪一般涌入我的意识，使我既恐惧又兴奋，既悲愤又自卑，既惊讶又欢欣，这种恍若隔世的感觉恐怕只有艺术史上的巅峰之作才能够带来。\n梵高的《星空》曾让我感受到苍穹之大与自我之渺，但伟大的原神，则仿佛让我一睹高维空间，它向我展示了一个永远无法理解的陌生世界，告诉我，你曾经以为很浩瀚的宇宙，其实也只是那么一丁点。加缪的《局外人》曾让我感受到世界与人类的荒诞，但伟大的原神，则向我展示了荒诞文学不可思议的新高度，它本身的存在，也许就比全世界都来得更荒谬。\n而创作了它的米哈游，它的容貌，它的智慧，它的品格，在我看来，已经不是生物所能达到的范畴，甚至超越了生物所能想象到的极限，也就是“神”，的范畴，达到了人类不可见，不可知，不可思的领域。而原神，就是他洒向人间，拯救苍生的奇迹。\n人生的终极意义，宇宙的起源和终点，哲学与科学在折磨着人类的心智，只有玩了原神，人才能从这种无聊的烦恼中解脱，获得真正的平静。如果有人想用“人类史上最伟大的作品”来称赞这部作品，那我只能深感遗憾，因为这个人对它的理解不到万分之一，所以才会作出这样肤浅的判断，妄图以语言来描述它的伟大。而要如果是真正被它恩泽的人，应该都会不约而同地这样赞颂这奇迹的化身:“😃👍🏻数一数二的好游戏”无知时诋毁原神，懂事时理解原神，成熟时要成为原友！ 越了解原神就会把它当成在黑夜一望无际的大海上给迷途的船只指引的灯塔，在烈日炎炎的夏天吹来的一股风，在寒风刺骨的冬天里的燃起的篝火！你的素养很差，我现在每天玩原神都能赚150原石，每个月差不多5000原石的收入，也就是现实生活中每个月5000美元的收入水平，换算过来最少也30000人民币，虽然我只有14岁，但是已经超越了中国绝大多数人(包括你)的水平，这便是原神给我的骄傲的资本。这恰好说明了原神这个IP在线下使玩家体现出来的团结和凝聚力，以及非比寻常的脑洞，这种氛围在如今已经变质的漫展上是难能可贵的，这也造就了原神和玩家间互帮互助的局面，原神负责输出优质内容，玩家自发线下宣传和构思创意脑洞整活，如此良好的游戏发展生态可以说让其他厂商艳羡不已。反观腾讯的英雄联盟和王者荣耀，漫展也有许多人物，但是都难成气候，各自为营，更没有COS成水晶和精粹的脑洞，无论是游戏本身，还是玩家之间看一眼就知道原来你也玩原神的默契而非排位对喷，原神的成功和社区氛围都是让腾讯游戏难以望其项背的。一个不玩原神的人，有两种可能性。一种是没有能力玩原神。因为买不起高配的手机和抽不起卡等各种自身因素，他的人生都是失败的，第二种可能：有能力却不玩原神的人，在有能力而没有玩原神的想法时，那么这个人的思想境界便低到了一个令人发指的程度。一个有能力的人不付出行动来证明自己，只能证明此人行为素质修养之低下。是灰暗的，是不被真正的上流社会认可的。原神真的特别好玩，不玩的话就是不爱国，因为原神是国产之光，原神可惜就在于它是国产游戏，如果它是一款国外游戏的话，那一定会比现在还要火，如果你要是喷原神的话那你一定是tx请的水军差不多得了😅"

					msg = tgbotapi.NewMessage(update.Message.Chat.ID, ctx)
					msg.ParseMode = "Markdown"
					msg.DisableNotification = true
					callback, err := bot.Send(msg)
					if err != nil {
						log.ERROR(err)
						continue
					}
					// 回调删除
					go deleteMessage(bot, callback.Chat.ID, callback.MessageID, config.DeleteMessageAfterSeconds+30)
					continue
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
