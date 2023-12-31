// Paimon_poi_test_bot
// https://core.telegram.org/bots/api#using-a-local-bot-api-server
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jellyqwq/Paimon/config"
	log "github.com/sirupsen/logrus"

	"github.com/jellyqwq/Paimon/news"
	"github.com/jellyqwq/Paimon/plugins/public/YearProgress"
	"github.com/jellyqwq/Paimon/plugins/self"
	"github.com/jellyqwq/Paimon/tools"
	"github.com/jellyqwq/Paimon/webapi"
	"github.com/robfig/cron"
)

var (
	// Control log month.
	LogMonth uint8

	// Regex of extracting inlineType and text.
	compileInlineInput = regexp.MustCompile(`^(?P<inlineType>.*?)(?: +(?P<text>.*)|$)`)
	compileElysia      = regexp.MustCompile(`^(派蒙|Paimon|飞行矮堇瓜|应急食品|白飞飞|神之嘴){1}`)

	// Inline keyboard about hotword.
	HotwordKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("B站", "HotWordBilibili"),
			tgbotapi.NewInlineKeyboardButtonData("微博", "HotWordWeibo"),
		),
	)

	// Inline keyboard about mihoyobbs exchange.
	HelpKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:                         "米游币兑换",
				SwitchInlineQueryCurrentChat: StringToPString("myb "),
			},
		),
	)

	//
	ypc = YearProgress.NewYearProgressConfig()
)

func init() {
	log.SetReportCaller(true)
	// Initalize log formatter.
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
		    //处理文件名
			fileName := path.Base(frame.File) + fmt.Sprintf(":%d", frame.Line)
			return frame.Function, fileName
		},
	})

	// Set log level.
	log.SetLevel(log.InfoLevel)
	logUpdate()
}

// update log
func logUpdate() {
	// Create directory logs if it does not exist.
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		if err := os.MkdirAll("./logs", 0666); err != nil {
			log.Warn(err)
		}
	}

	// Create log files if it does not exist, which is named yyyy-MM and the suffix is log.
	nowTime := time.Now().Local()
	LogMonth = uint8(nowTime.Month())
	if file, err := os.OpenFile(
		fmt.Sprintf("./logs/%s.log", nowTime.Format("2006-01")),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	); err == nil {
		writers := []io.Writer{
			file,
			os.Stdout,
		}
		log.SetOutput(io.MultiWriter(writers...))
	} else {
		log.Warn("Failed to log to file, using default stderr")
	}
}

func StringToPString(s string) *string {
	return &s
}

// Delete telegram message after delay sceonds.
func deleteMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int, delay int64) {
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)
	time.Sleep(time.Duration(delay) * time.Second)
	bot.Send(msg)
}

// func InitMessage(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
// 	msg.ParseMode = "Markdown"
// 	msg.DisableWebPagePreview = true
// 	msg.DisableNotification = true
// 	return msg
// }

// Main Handler
func mainHandler() {
	// Read the config.
	config, err := config.ReadYaml()
	if err != nil {
		log.Fatal(err)
	}

	// Update mihoyo bbs goods information.
	// if err = webapi.MihoyoBBSGoodsUpdate(); err != nil {
	// 	log.Fatal(err)
	// }

	// Create a variable named bot of *tgbotapi.BotAPI.
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	// Control debugging information of bot.
	bot.Debug = true
	log.Infof("Authorized on account %s", bot.Self.UserName)

	// webhook, _ := tgbotapi.NewWebhookWithCert(config.Webhook.URL + bot.Token, tgbotapi.FilePath(config.Webhook.CertificatePemPath))
	webhook, _ := tgbotapi.NewWebhook(config.TelegramWebHook.Url + bot.Token)
	webhook.IPAddress = config.TelegramWebHook.IPAddress
	webhook.AllowedUpdates = config.TelegramWebHook.AllowedUpdates
	webhook.MaxConnections = config.TelegramWebHook.MaxConnections

	// Request sends a Chattable to Telegram.
	if _, err = bot.Request(webhook); err != nil {
		log.Fatal(err)
	}

	// Create a timer.
	c := cron.New()

	// Add a task for updating mihoyobbs goods information.
	if err = c.AddFunc("0 0 4 * * ?", func() {
		if err = webapi.MihoyoBBSGoodsUpdate(); err != nil {
			log.Error(err)
		} else {
			log.Info("Cron MihoyoBBSGoodsUpdate is sent.")
		}
	}); err != nil {
		log.Error("AddFunc error : ", err)
	} else {
		log.Info("Cron MihoyoBBSGoodsUpdate is loaded.")
	}

	// Add a task for acquiring string of YearProgressBar.
	ypc.ChatID = config.YearProgressChatId
	if err = c.AddFunc("0 */2 * * * ?", func() {
		if bar := ypc.GetYearProgress(); bar != "" {
			bar = fmt.Sprintf("*YearProgress*\n%s", bar)
			msg := tgbotapi.NewMessage(ypc.ChatID, bar)
			msg.ParseMode = "Markdown"
			msg.DisableNotification = true
			if _, err := bot.Send(msg); err != nil {
				log.Error("Send error : ", err)
			} else {
				log.Info("Cron GetYearProgress is sent.")
			}
		}
	}); err != nil {
		log.Error("AddFunc error : ", err)
	} else {
		log.Info("Cron ypc.GetYearProgress is loaded.")
	}

	c.Start()
	defer c.Stop()

	// Set route of webhook.
	updates := bot.ListenForWebhook("/telegram/" + bot.Token)
	go http.ListenAndServe(config.WebhookIP+":"+strconv.FormatUint(config.WebhookPort, 10), nil)

	// Receive updates and process them.
	for update := range updates {
		// Update log if the variable LogMonth is not equal current month.
		if uint8(time.Now().Local().Month()) != LogMonth {
			logUpdate()
		}

		// Process message if it is not null.
		if update.Message != nil {
			// Inline keyboard with command.
			if update.Message.IsCommand() {
				log.Debug(update.Message.Command())

				// Classify command handle.
				switch update.Message.Command() {

				// Return nas ip, this is a plugin.
				case "nas":
					{
						// Firstly delete command after seconds.
						go deleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID, config.DeleteMessageAfterSeconds)

						// Get my nas ipv6
						text := self.GetNASIpv6(update.Message.Chat.ID, update.Message.From.ID)
						if text == "" {
							continue
						}

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
						msg.ParseMode = "Markdown"
						msg.DisableNotification = true

						rep, err := bot.Send(msg)
						if err != nil {
							log.Error(err)
							continue
						}

						go deleteMessage(bot, rep.Chat.ID, rep.MessageID, config.DeleteMessageAfterSeconds)
					}

				// Return InlineKeyboard that contains hotword of bilibili and weibo.
				case "hot_word":
					{
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "热搜🔥 | 大瓜🍉")
						msg.ReplyMarkup = HotwordKeyboard
						msg.DisableNotification = true

						rep, err := bot.Send(msg)
						if err != nil {
							log.Error(err)
							continue
						}

						go deleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID, config.DeleteMessageAfterSeconds)
						go deleteMessage(bot, rep.Chat.ID, rep.MessageID, config.DeleteMessageAfterSeconds)
					}

				// Return common currency exchange rates.
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
							log.Error(err)
							continue
						}

						go deleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID, config.DeleteMessageAfterSeconds)
						go deleteMessage(bot, rep.Chat.ID, rep.MessageID, config.DeleteMessageAfterSeconds)
					}

				// Return help menu.
				case "help":
					{
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "[Paimon | 应急食品](https://github.com/jellyqwq/Paimon)")
						msg.ParseMode = "Markdown"
						msg.ReplyMarkup = HelpKeyboard
						msg.DisableWebPagePreview = true
						msg.DisableNotification = true

						rep, err := bot.Send(msg)
						if err != nil {
							log.Error(err)
							continue
						}

						go deleteMessage(bot, update.Message.Chat.ID, update.Message.MessageID, config.DeleteMessageAfterSeconds)
						go deleteMessage(bot, rep.Chat.ID, rep.MessageID, config.DeleteMessageAfterSeconds)
					}
				}

				// Terminate this processing.
				continue
			}

			// Process message text which is not empty.
			if update.Message.Text != "" {
				text := update.Message.Text

				// Replace bot keyword.
				text = string(compileElysia.ReplaceAll([]byte(text), []byte("")))
				log.Debug(fmt.Sprintf(`{"chat_id": "%v", "name": "%v", "message": "%v"}`, update.Message.Chat.ID, update.Message.From.FirstName, text))

				var msg tgbotapi.MessageConfig

				// Return telegram user information to sender or responder.
				if strings.Contains(text, "INFO") {
					var ctx string = ""
					if update.Message.ReplyToMessage != nil {
						ctx = fmt.Sprintf("ReplyUserInfo\nUserID:`%v`\nChatID:`%v`\nFirstName:`%v`\nLastName:`%v`\nUserName:`%v`", update.Message.ReplyToMessage.From.ID, update.Message.ReplyToMessage.Chat.ID, update.Message.ReplyToMessage.From.FirstName, update.Message.ReplyToMessage.From.LastName, update.Message.ReplyToMessage.From.UserName)
					} else {
						ctx = fmt.Sprintf("UserInfo\nUserID:`%v`\nChatID:`%v`\nFirstName:`%v`\nLastName:`%v`\nUserName:`%v`", update.Message.From.ID, update.Message.Chat.ID, update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName)
					}

					msg = tgbotapi.NewMessage(update.Message.Chat.ID, ctx)
					msg.ParseMode = "Markdown"
					msg.DisableNotification = true
					msg.ReplyToMessageID = update.Message.MessageID
					if msg.Text != "" {
						if callback, err := bot.Send(msg); err != nil {
							log.Error(err)
							continue
						} else {
							go deleteMessage(bot, callback.Chat.ID, callback.MessageID, config.DeleteMessageAfterSeconds+30)
							continue
						}
					}

					// Terminate this processing.
					continue
				}

				// Return the OP message if sender message contains only the word "原神".
				if text == "原神" {
					ctx := "差不多得了😅屁大点事都要拐上原神，原神一没招你惹你，二没干伤天害理的事情，到底怎么你了让你一直无脑抹黑，米哈游每天费尽心思的文化输出弘扬中国文化，你这种喷子只会在网上敲键盘诋毁良心公司，中国游戏的未来就是被你这种人毁掉的😅\n叫我们原批的小心点\n老子在大街上亲手给打过两个\n我在公共座椅上无聊玩原神，有两个B就从我旁边过，看见我玩原神就悄悄说了一句:又是一个原批，我就直接上去一拳呼脸上，我根本不给他解释的机会，我也不问他为什么说我是原批，我就打，我就看他不爽，他惹我了，我就不给他解释的机会，直接照着脸和脑门就打直接给那B呼出鼻血，脸上青一块，紫一块的我没撕她嘴巴都算好了你们这还不算最狠的，我记得我以前小时候春节去老家里，有一颗核弹，我以为是鞭炮，和大地红一起点了，当时噼里啪啦得，然后突然一朵蘑菇云平地而起，当时我就只记得两眼一黑，昏过去了，整个村子没了，幸好我是体育生，身体素质不错，住了几天院就没事了，几个月下来腿脚才利落，现在已经没事了，但是那种钻心的疼还是让我一生难忘，😂😂😂  令人感叹今早一玩原神，我便昏死了过去，现在才刚刚缓过来。在昏死过去的短短数小时内，我的大脑仿佛被龙卷风无数次摧毁。\n在原神这一神作的面前，我就像一个一丝不挂的原始人突然来到了现代都市，二次元已如高楼大厦将我牢牢地吸引，开放世界就突然变成那喇叭轰鸣的汽车，不仅把我吓个措手不及，还让我瞬间将注意完全放在了这新的奇物上面，而还没等我稍微平复心情，纹化输出的出现就如同眼前遮天蔽日的宇宙战舰，将我的世界观无情地粉碎，使我彻底陷入了忘我的迷乱，狂泄不止。\n原神，那眼花缭乱的一切都让我感到震撼，但是我那贫瘠的大脑却根本无法理清其中任何的逻辑，巨量的信息和情感泄洪一般涌入我的意识，使我既恐惧又兴奋，既悲愤又自卑，既惊讶又欢欣，这种恍若隔世的感觉恐怕只有艺术史上的巅峰之作才能够带来。\n梵高的《星空》曾让我感受到苍穹之大与自我之渺，但伟大的原神，则仿佛让我一睹高维空间，它向我展示了一个永远无法理解的陌生世界，告诉我，你曾经以为很浩瀚的宇宙，其实也只是那么一丁点。加缪的《局外人》曾让我感受到世界与人类的荒诞，但伟大的原神，则向我展示了荒诞文学不可思议的新高度，它本身的存在，也许就比全世界都来得更荒谬。\n而创作了它的米哈游，它的容貌，它的智慧，它的品格，在我看来，已经不是生物所能达到的范畴，甚至超越了生物所能想象到的极限，也就是“神”，的范畴，达到了人类不可见，不可知，不可思的领域。而原神，就是他洒向人间，拯救苍生的奇迹。\n人生的终极意义，宇宙的起源和终点，哲学与科学在折磨着人类的心智，只有玩了原神，人才能从这种无聊的烦恼中解脱，获得真正的平静。如果有人想用“人类史上最伟大的作品”来称赞这部作品，那我只能深感遗憾，因为这个人对它的理解不到万分之一，所以才会作出这样肤浅的判断，妄图以语言来描述它的伟大。而要如果是真正被它恩泽的人，应该都会不约而同地这样赞颂这奇迹的化身:“😃👍🏻数一数二的好游戏”无知时诋毁原神，懂事时理解原神，成熟时要成为原友！ 越了解原神就会把它当成在黑夜一望无际的大海上给迷途的船只指引的灯塔，在烈日炎炎的夏天吹来的一股风，在寒风刺骨的冬天里的燃起的篝火！你的素养很差，我现在每天玩原神都能赚150原石，每个月差不多5000原石的收入，也就是现实生活中每个月5000美元的收入水平，换算过来最少也30000人民币，虽然我只有14岁，但是已经超越了中国绝大多数人(包括你)的水平，这便是原神给我的骄傲的资本。这恰好说明了原神这个IP在线下使玩家体现出来的团结和凝聚力，以及非比寻常的脑洞，这种氛围在如今已经变质的漫展上是难能可贵的，这也造就了原神和玩家间互帮互助的局面，原神负责输出优质内容，玩家自发线下宣传和构思创意脑洞整活，如此良好的游戏发展生态可以说让其他厂商艳羡不已。反观腾讯的英雄联盟和王者荣耀，漫展也有许多人物，但是都难成气候，各自为营，更没有COS成水晶和精粹的脑洞，无论是游戏本身，还是玩家之间看一眼就知道原来你也玩原神的默契而非排位对喷，原神的成功和社区氛围都是让腾讯游戏难以望其项背的。一个不玩原神的人，有两种可能性。一种是没有能力玩原神。因为买不起高配的手机和抽不起卡等各种自身因素，他的人生都是失败的，第二种可能：有能力却不玩原神的人，在有能力而没有玩原神的想法时，那么这个人的思想境界便低到了一个令人发指的程度。一个有能力的人不付出行动来证明自己，只能证明此人行为素质修养之低下。是灰暗的，是不被真正的上流社会认可的。原神真的特别好玩，不玩的话就是不爱国，因为原神是国产之光，原神可惜就在于它是国产游戏，如果它是一款国外游戏的话，那一定会比现在还要火，如果你要是喷原神的话那你一定是tx请的水军差不多得了😅"

					msg = tgbotapi.NewMessage(update.Message.Chat.ID, ctx)
					msg.ParseMode = "Markdown"
					msg.DisableNotification = true
					callback, err := bot.Send(msg)
					if err != nil {
						log.Error(err)
						continue
					}

					// Delete the OP message shortly after.
					go deleteMessage(bot, callback.Chat.ID, callback.MessageID, config.DeleteMessageAfterSeconds+30)

					// Terminate this processing
					continue
				}

				// Terminate this processing
				continue
			}

			// Terminate this processing.
			continue
		}

		// Process CallbackQuery if it is not null.
		if update.CallbackQuery != nil {
			// Set callback loading information.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "loading...")
			if _, err := bot.Request(callback); err != nil {
				log.Error(err)
				continue
			}

			// Classify callback application.
			CallbackQueryData := update.CallbackQuery.Data
			if CallbackQueryData == "HotWordBilibili" {
				ctx, err := news.BiliHotWords()
				if err != nil {
					log.Error(err)
					continue
				}

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx)
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				msg.DisableNotification = true
				if msg.Text != "" {
					if _, err := bot.Send(msg); err != nil {
						log.Error(err)
						continue
					}
				}

				// Terminate this processing.
				continue
			}

			if CallbackQueryData == "HotWordWeibo" {
				ctx, err := news.WeiboHotWords()
				if err != nil {
					log.Error(err)
					continue
				}

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx)
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				msg.DisableNotification = true
				if msg.Text != "" {
					if _, err := bot.Send(msg); err != nil {
						log.Error(err)
						continue
					}
				}

				// Terminate this processing.
				continue
			}

			if len(CallbackQueryData) > 10 {
				// Process currency query.
				if CallbackQueryData[:8] == "currency" {
					tempList := strings.Split(CallbackQueryData, "-")
					currency := tempList[1] + "-" + tempList[2]
					ctx, err := webapi.Finance(currency)
					if err != nil {
						log.Error(err)
						continue
					}
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ctx)
					msg.ParseMode = "Markdown"
					msg.DisableWebPagePreview = true
					msg.DisableNotification = true
					if msg.Text != "" {
						if _, err := bot.Send(msg); err != nil {
							log.Error(err)
							continue
						}
					}

					// Terminate this processing
					continue
				}

				// Terminate this processing
				continue
			}

			// Terminate this processing.
			continue
		}

		// Process InlineQuery if it is not null.
		if update.InlineQuery != nil {
			text := update.InlineQuery.Query

			if text == "" {
				continue
			}

			// Extract type and text from the query.
			paramsMap := tools.GetParamsOneDimension(compileInlineInput, text)
			inlineType := paramsMap["inlineType"]
			text = paramsMap["text"]
			if inlineType == "" {
				continue
			}

			log.Debug("inlineType", inlineType)
			log.Debug("text", text)

			result := []interface{}{}
			switch inlineType {
			case "myb":
				{
					result, err = webapi.MihoyoBBSGoodsForQuery(text)
					if err != nil {
						log.Error(err)
						continue
					}
				}
			}

			if len(result) == 0 {
				log.Error("result is empty")
				continue
			}

			// Structure a InlineConfig.
			inlineConf := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     10,
				Results:       result,
			}
			if _, err := bot.Send(inlineConf); err != nil {
				log.Error(err)
				continue
			}

			// Terminate this processing.
			continue
		}
	}
}

func main() {
	mainHandler()
}
