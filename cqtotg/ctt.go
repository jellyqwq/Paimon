package cqtotg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"

	// "strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	GroupID   int `json:"group_id"`
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

type PostParams struct {
	Bot  *tgbotapi.BotAPI
	Conf config.Config
}

type MessageOutput struct {
	ImageSlice []string `json:"image_slice"`
	GIFSlice   []string `json:"gif_slice"`
	Text       string   `json:"text"`
}

type Notice struct {
	PostType   string `json:"post_type"`
	NoticeType string `json:"notice_type"`
	Time       int    `json:"time"`
	SelfID     int64  `json:"self_id"`
	GroupID    int    `json:"group_id"`
	UserID     int    `json:"user_id"`
	File       struct {
		Busid int    `json:"busid"`
		ID    string `json:"id"`
		Name  string `json:"name"`
		Size  int    `json:"size"`
		URL   string `json:"url"`
	} `json:"file"`
}

var File2ContentType = map[string]string{
	"mp4":  "video/mp4",
	"png":  "image/png",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
}

func (bot *PostParams) Post(writer http.ResponseWriter, request *http.Request) {
	x, _ := io.ReadAll(request.Body)

	var jsonRet PublicHeartbeatRequest
	json.Unmarshal(x, &jsonRet)

	PostType := jsonRet.PostType

	switch PostType {

	case "message":
		{
			var output MessageOutput

			var message UnheartbeatRequest
			json.Unmarshal(x, &message)

			// replace photo and add photo url to output.ImageSlice
			compileImage := regexp.MustCompile(`\[CQ:image,file=[0-9a-z]+\.image,(subType=[0-9]+,)?url=(?P<url>https:\/\/(c2cpicdw|gchat).qpic.cn\/.*?)\]`)
			str := compileImage.ReplaceAllStringFunc(message.Message, func(s string) string {
				paramsMap := tools.GetParamsOneDimension(compileImage, s)
				url := paramsMap["url"]

				// classify image type to output.ImageSlice or output.GIFSlice by http response.Header
				response, _ := http.Get(url)

				typeSlice := response.Header["Content-Type"]
				if tools.IsOneDimensionSliceContainsString(typeSlice, "image/gif") {
					output.GIFSlice = append(output.GIFSlice, url)
				} else {
					output.ImageSlice = append(output.ImageSlice, url)
				}
				return ""
			})

			// replace mention message
			compileMention := regexp.MustCompile(`\[CQ:at,qq=(?P<qq>[0-9]+)\]`)
			str = compileMention.ReplaceAllStringFunc(str, func(s string) string {
				paramsMap := tools.GetParamsOneDimension(compileMention, s)
				i, _ := strconv.Atoi(paramsMap["qq"])
				if _, ok := bot.Conf.CQ2TG.MentionReflection[uint64(i)]; ok {
					return "[靓仔](tg://user?id=" + strconv.FormatUint(bot.Conf.CQ2TG.MentionReflection[uint64(i)], 10) + ")"
				} else {
					return ""
				}
			})

			GIFSliceLength := len(output.GIFSlice)
			ImageSliceLength := len(output.ImageSlice)

			faceDeleted := regexp.MustCompile(`\[CQ:face,id=[0-9]+\]`)
			str = string(faceDeleted.ReplaceAll([]byte(str), []byte("")))

			gstr := ""
			if message.GroupID != 0 {
				gstr = fmt.Sprintf("`%v` | ", message.GroupID)
			}
			output.Text = fmt.Sprintf("Forward from %v*%v* `%v`\n%v", gstr, message.Sender.Nickname, message.UserID, str)
            if str == "" {
                output.Text = ""
            }

			var ImageList []interface{}

			// Image message send
			if ImageSliceLength == 1 {
				msg := tgbotapi.NewPhoto(bot.Conf.CQ2TG.RecivedChatId, tgbotapi.FileURL(output.ImageSlice[0]))
				msg.Caption = output.Text
				msg.DisableNotification = true
				msg.ParseMode = "Markdown"
				bot.Bot.Send(msg)

			} else if ImageSliceLength > 1 && ImageSliceLength <= 10 {
				captionOK := false
				for _, value := range output.ImageSlice {
					Image := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(value))
					if !captionOK {
						Image.Caption = output.Text
						Image.ParseMode = "Markdown"
						captionOK = true
					}
					ImageList = append(ImageList, Image)
				}
				msg := tgbotapi.NewMediaGroup(bot.Conf.CQ2TG.RecivedChatId, ImageList)
				msg.DisableNotification = true
				bot.Bot.Send(msg)
			}

			if GIFSliceLength == 1 {
				msg := tgbotapi.NewDocument(bot.Conf.CQ2TG.RecivedChatId, tgbotapi.FileURL(output.GIFSlice[0]))
				msg.Caption = output.Text
				msg.DisableNotification = true
				msg.ParseMode = "Markdown"
				bot.Bot.Send(msg)
			} else if 1 < GIFSliceLength && GIFSliceLength <= 10 {
				captionOK := false
				for _, value := range output.GIFSlice {
					GIF := tgbotapi.NewInputMediaDocument(tgbotapi.FileURL(value))
					if !captionOK {
						GIF.Caption = output.Text
						GIF.ParseMode = "Markdown"
						captionOK = true
					}
					ImageList = append(ImageList, GIF)
				}
				msg := tgbotapi.NewMediaGroup(bot.Conf.CQ2TG.RecivedChatId, ImageList)
				msg.DisableNotification = true
				bot.Bot.Send(msg)
			}

			if GIFSliceLength == 0 && ImageSliceLength == 0 {
				msg := tgbotapi.NewMessage(bot.Conf.CQ2TG.RecivedChatId, output.Text)
				msg.DisableNotification = true
				msg.ParseMode = "Markdown"
				if output.Text != "" {
					bot.Bot.Send(msg)
				}
			}
		}

	case "notice":
		{
			var notice Notice

			json.Unmarshal(x, &notice)

			if notice.File.Size <= 52428800 {
				// mp4, png, jpg, zip, ...都是notice, 所以在这里要加一个正则表达式来分类
				// filename := notice.File.Name

				// compileFileFormat := regexp.MustCompile(`.+\.(?P<format>.+)$`)
				// paramsMap := tools.GetParamsOneDimension(compileFileFormat, filename)
				// fileformat := paramsMap["format"]
				// fileformat = strings.ToLower(fileformat)
				resp, err := http.Get(notice.File.URL)
				if err != nil {
					log.Panicln(err)
				}
				defer resp.Body.Close()
				rb, _ := io.ReadAll(resp.Body)

				file := tgbotapi.FileReader{
					Name:   notice.File.Name,
					Reader: bytes.NewReader(rb),
				}

				// msg := tgbotapi.NewDocument(bot.Conf.CQ2TG.RecivedChatId, tgbotapi.FileURL(bot.Conf.TelegramWebHook.Url + "retype/?url=" + notice.File.URL + "&filename=" + notice.File.Name + "&type=" + value))
				msg := tgbotapi.NewDocument(bot.Conf.CQ2TG.RecivedChatId, file)
				gstr := ""
				if notice.GroupID != 0 {
					gstr = fmt.Sprintf("`%v` | ", notice.GroupID)
				}
				msg.Caption = fmt.Sprintf("Forward `%v` from %v`%v`", notice.File.Name, gstr, notice.UserID)
				msg.DisableNotification = true
				msg.ParseMode = "Markdown"
				bot.Bot.Send(msg)
				// if value, ok := File2ContentType[fileformat]; ok {
				//     if fileformat == "mp4" {
				//         msg := tgbotapi.NewVideo(bot.Conf.CQ2TG.RecivedChatId, tgbotapi.FileURL(bot.Conf.TelegramWebHook.Url + "retype/?url=" + notice.File.URL + "&filename=" + notice.File.Name + "&type=" + value))
				//         gstr := ""
				//         if notice.GroupID != 0 {
				//             gstr = fmt.Sprintf("`%v` | ", notice.GroupID)
				//         }
				//         msg.Caption = fmt.Sprintf("Forward `%v` from %v`%v`", notice.File.Name, gstr, notice.UserID)
				//         msg.DisableNotification = true
				//         msg.ParseMode = "Markdown"
				//         msg.SupportsStreaming = true
				//         bot.Bot.Send(msg)
				//     } else {
				//         resp, err := http.Get(notice.File.URL)
				//         if err != nil {
				//             log.Panicln(err)
				//         }
				//         defer resp.Body.Close()
				//         rb, _ := io.ReadAll(resp.Body)

				//         file := tgbotapi.FileReader {
				//             Name: notice.File.Name,
				//             Reader: bytes.NewReader(rb),
				//         }

				//         // msg := tgbotapi.NewDocument(bot.Conf.CQ2TG.RecivedChatId, tgbotapi.FileURL(bot.Conf.TelegramWebHook.Url + "retype/?url=" + notice.File.URL + "&filename=" + notice.File.Name + "&type=" + value))
				//         msg := tgbotapi.NewDocument(bot.Conf.CQ2TG.RecivedChatId, file)
				//         gstr := ""
				//         if notice.GroupID != 0 {
				//             gstr = fmt.Sprintf("`%v` | ", notice.GroupID)
				//         }
				//         msg.Caption = fmt.Sprintf("Forward `%v` from %v`%v`", notice.File.Name, gstr, notice.UserID)
				//         msg.DisableNotification = true
				//         msg.ParseMode = "Markdown"
				//         bot.Bot.Send(msg)
				//     }
				// }

			}
		}
	}
}

// reset video response Header
func FileParse(writer http.ResponseWriter, request *http.Request) {
	params := request.URL.Query()
	url := params["url"][0]
	fileName := params["filename"][0]
	fileType := params["type"][0]

	resp, err := http.Get(url)
	if err != nil {
		log.Panicln(err)
	}

	defer resp.Body.Close()

	writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	writer.Header().Set("Content-Type", fileType)
	writer.Header().Set("Content-Length", resp.Header.Values("Content-Length")[0])
	writer.Header().Set("Connection", "keep-alive")
	rb, _ := io.ReadAll(resp.Body)
	_, err = writer.Write(rb)
	if err != nil {
		log.Panicln(err)
	}
}
