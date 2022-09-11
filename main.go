// Paimon_poi_test_bot
// https://core.telegram.org/bots/api#using-a-local-bot-api-server
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	news "github.com/jellyqwq/Paimon/news"
)

type config struct {
	BotToken string `json:"BotToken"`
	Webhook  struct {
		URL            string   `json:"url"`
		CertificatePemPath    string   `json:"certificatePemPath"`
		CertificateKeyPath string `json:"certificateKeyPath"`
		IPAddress      string   `json:"ip_address"`
		MaxConnections int      `json:"max_connections"`
		AllowedUpdates []string `json:"allowed_updates"`
	} `json:"webhook"`
}

func logError(err error) () {
	if err != nil {
	  	log.Fatal(err)
	}
}

func readConfig() (*config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	
	conf := config{}
	err = decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func main() {
	config, err := readConfig()
	logError(err)

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	logError(err)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// webhook, _ := tgbotapi.NewWebhookWithCert(config.Webhook.URL + bot.Token, tgbotapi.FilePath(config.Webhook.CertificatePemPath))
	webhook, _ := tgbotapi.NewWebhook(config.Webhook.URL + bot.Token)
	// webhook.IPAddress = config.Webhook.IPAddress
	webhook.AllowedUpdates = config.Webhook.AllowedUpdates
	webhook.MaxConnections = config.Webhook.MaxConnections

	_, err = bot.Request(webhook)
	logError(err)

	info, err := bot.GetWebhookInfo()
	logError(err)

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	// cqhttp http-reverse
	http.HandleFunc("/cq/", Post)

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe("127.0.0.1:6700", nil)

	for update := range updates {
		log.Printf("%v", update)
		if update.Message != nil {
			text := update.Message.Text

			regElysia := regexp.MustCompile(`^(爱|e|E){1}(莉|ly){1}(希雅|sia)?`)
			if (regElysia.Match([]byte(text))) {
				text = string(regElysia.ReplaceAll([]byte(text), []byte("")))

				var msg tgbotapi.MessageConfig
				if strings.Contains(text, "bhot") || strings.Contains(text, "鼠鼠热搜") {
					ctx, err:= news.BiliHotWords()
					logError(err)

					msg = tgbotapi.NewMessage(update.Message.Chat.ID, ctx)
					msg.ParseMode = "Markdown"
					msg.DisableWebPagePreview = true

				} else if strings.Contains(text, "whot") || strings.Contains(text, "微博"){
					ctx, err := news.WeiboHotWords()
					logError(err)

					msg = tgbotapi.NewMessage(update.Message.Chat.ID, ctx)
					msg.ParseMode = "Markdown"
					msg.DisableWebPagePreview = true

				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "你好,我是爱莉希雅")
				}
				_, err = bot.Send(msg)
				logError(err)
			}
		}
	}
}

func Post(writer http.ResponseWriter, request *http.Request) {

	// log.Printf("=========================================")
	// log.Printf("nn%v\n", request)
	// log.Printf("-----------------------------------------")
	x, _ := io.ReadAll(request.Body)
	jsonRet := map[string]interface{}{}
	json.Unmarshal(x, &jsonRet)
	log.Printf("%T: %v\n",jsonRet ,jsonRet)
	log.Println(jsonRet["message"])
	log.Printf("=========================================")

}