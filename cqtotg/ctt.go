package cqtotg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jellyqwq/Paimon/config"
	"github.com/jellyqwq/Paimon/tools"
)

type PublicHeartbeatRequest struct {
    PostType    string `json:"post_type"`
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

type PostParams struct {
    Bot *tgbotapi.BotAPI
    Conf config.Config
}

type MessageOutput struct {
    ImageSlice []string `json:"image_slice"`
    GIFSlice []string `json:"gif_slice"`
    Text string `json:"text"`

}

func (bot *PostParams) Post(writer http.ResponseWriter, request *http.Request) {
    x, _ := io.ReadAll(request.Body)
    
    var jsonRet PublicHeartbeatRequest
    json.Unmarshal(x, &jsonRet)

    PostType := jsonRet.PostType

    switch PostType {

        case "message": {
            var output MessageOutput

            var message UnheartbeatRequest
            json.Unmarshal(x, &message)

            // replace photo and add photo url to output.ImageSlice
            compileImage := regexp.MustCompile(`\[CQ:image,file=[0-9a-z]+\.image,(subType=[0-9]+,)?url=(?P<url>https:\/\/(c2cpicdw|gchat).qpic.cn\/.*?)\]`)
            str := compileImage.ReplaceAllStringFunc(message.Message, func(s string) string {
                paramsMap := tools.GetParamsOneDimension(compileImage, s)
                url := paramsMap["url"]

                // classify image type to output.ImageSlice or output.GIFSlice by http response.Header
                response, err := http.Get(url)
                if err != nil {
                    log.Println(err)
                }

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
            output.Text = fmt.Sprintf("`%v` %v sent:\n%v", message.UserID, message.Sender.Nickname, str)
            
            var ImageList []interface{}

            // Image message send
            if ImageSliceLength == 0 {
                msg := tgbotapi.NewMessage(bot.Conf.CQ2TG.RecivedChatId, output.Text)
                msg.DisableNotification = true
                msg.ParseMode = "Markdown"
                bot.Bot.Send(msg)

            } else if ImageSliceLength == 1 {
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

            if GIFSliceLength == 0 {
                msg := tgbotapi.NewMessage(bot.Conf.CQ2TG.RecivedChatId, output.Text)
                msg.DisableNotification = true
                msg.ParseMode = "Markdown"
                bot.Bot.Send(msg)
            } else if GIFSliceLength == 1 {
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
        }
    }
}
