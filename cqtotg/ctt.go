package cqtotg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/jellyqwq/Paimon/config"
	"github.com/jellyqwq/Paimon/tools"
)

type PublicHeartbeatRequest struct {
	PostType string `json:"post_type"`
}

type UnheartbeatRequest struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	Time        int    `json:"time"`
	SelfID      int64  `json:"self_id"`
	SubType     string `json:"sub_type"`
	TargetID    int64  `json:"target_id"`
	Message     string `json:"message"`
	RawMessage  string `json:"raw_message"`
	Font        int    `json:"font"`
	Sender      struct {
		Age      int    `json:"age"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		UserID   int    `json:"user_id"`
	} `json:"sender"`
	MessageID int `json:"message_id"`
	UserID    int `json:"user_id"`
}

type HeartbeatRequest struct {
	PostType      string `json:"post_type"`
	MetaEventType string `json:"meta_event_type"`
	Time          int    `json:"time"`
	SelfID        int64  `json:"self_id"`
	Interval      int    `json:"interval"`
	Status        struct {
		AppEnabled     bool        `json:"app_enabled"`
		AppGood        bool        `json:"app_good"`
		AppInitialized bool        `json:"app_initialized"`
		Good           bool        `json:"good"`
		Online         bool        `json:"online"`
		PluginsGood    interface{} `json:"plugins_good"`
		Stat           struct {
			PacketReceived  int `json:"packet_received"`
			PacketSent      int `json:"packet_sent"`
			PacketLost      int `json:"packet_lost"`
			MessageReceived int `json:"message_received"`
			MessageSent     int `json:"message_sent"`
			LastMessageTime int `json:"last_message_time"`
			DisconnectTimes int `json:"disconnect_times"`
			LostTimes       int `json:"lost_times"`
		} `json:"stat"`
	} `json:"status"`
}

type TelegramBot struct {
	Bot  *tgbotapi.BotAPI
	Conf config.Config
}

func (bot *TelegramBot) Post(writer http.ResponseWriter, request *http.Request) {
	var jsonRet PublicHeartbeatRequest
	var message UnheartbeatRequest

	// make sure the body won't be closed until itself exits
	defer request.Body.Close()
	for {
		x, err := io.ReadAll(request.Body)
		if err != nil {
			return
		}

		// Skip those invaild requests
		if err := json.Unmarshal(x, &jsonRet); err != nil {
			fmt.Println(err)
			continue
		}

		PostType := jsonRet.PostType

		switch PostType {

		case "message":
			{
				// log.Printf("%v", string(x))

				json.Unmarshal(x, &message)

				// 这是单独对图片进行匹配的
				reg := regexp.MustCompile(`\[CQ:(?P<type>image),file=[0-9a-z]+\.image,(subType=[0-9]+,)?url=(?P<url>https:\/\/(c2cpicdw|gchat).qpic.cn\/.*?)\]`)
				str := reg.ReplaceAllString(message.Message, "")

				// 对文本
				// regAt := regexp.MustCompile(`\[CQ:at,qq=([0-9]+)\]`)
				regAt := regexp.MustCompile(`\[CQ:at,qq=(?P<qq>[0-9]+)\]`)

				str = regAt.ReplaceAllStringFunc(str, func(s string) string {
					paramsMap := tools.GetParamsOneDimension(regAt, s)
					i, _ := strconv.Atoi(paramsMap["qq"])
					fmt.Println(i)
					if _, ok := bot.Conf.CQ2TG.MentionReflection[uint64(i)]; ok {
						return "[靓仔](tg://user?id=" + strconv.FormatUint(bot.Conf.CQ2TG.MentionReflection[uint64(i)], 10) + ")"
					} else {
						return ""
					}

				})

				ctx := fmt.Sprintf("`%v` %v said:\n%v", message.UserID, message.Sender.Nickname, str)

				groupsName := reg.SubexpNames()

				var IamgeSlice []string
				// var GIFSlice []string

				for _, match := range reg.FindAllStringSubmatch(message.Message, -1) {
					for groupIndex, value := range match {
						// 该循环用于判断CQ类型并处理
						key := groupsName[groupIndex]
						// captionOK := false

						switch key {
						case "type":
							{
								switch value {
								case "image":
									continue
								default:
									break
								}
							}
						case "url":
							{
								IamgeSlice = append(IamgeSlice, value)
							}
						}
					}
				}

				IamgeSliceLength := len(IamgeSlice)
				fmt.Println(IamgeSlice, IamgeSliceLength)

				// var ImagesListContainer [][]interface{}
				var ImagesList []interface{}

				if IamgeSliceLength == 0 {
					msg := tgbotapi.NewMessage(bot.Conf.CQ2TG.RecivedChatId, ctx)
					msg.DisableNotification = true
					msg.ParseMode = "Markdown"
					bot.Bot.Send(msg)

				} else if IamgeSliceLength == 1 {
					msg := tgbotapi.NewPhoto(bot.Conf.CQ2TG.RecivedChatId, tgbotapi.FileURL(IamgeSlice[0]))
					msg.Caption = ctx
					msg.DisableNotification = true
					msg.ParseMode = "Markdown"
					bot.Bot.Send(msg)

				} else if IamgeSliceLength > 1 && IamgeSliceLength <= 10 {
					captionOK := false
					for _, value := range IamgeSlice {
						Image := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(value))
						if !captionOK {
							Image.Caption = ctx
							Image.ParseMode = "Markdown"
							captionOK = true
						}
						ImagesList = append(ImagesList, Image)
					}
					msg := tgbotapi.NewMediaGroup(bot.Conf.CQ2TG.RecivedChatId, ImagesList)
					msg.DisableNotification = true
					bot.Bot.Send(msg)
				}

			}
		}
	}
}
