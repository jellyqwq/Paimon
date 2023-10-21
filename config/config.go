package config

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BotToken        string `yaml:"BotToken"`
	WebhookIP       string `yaml:"WebhookIP"`
	WebhookPort     uint64 `yaml:"WebhookPort"`
	TelegramWebHook struct {
		Url            string   `yaml:"Url"`
		IPAddress      string   `yaml:"IPAddress"`
		MaxConnections int      `yaml:"MaxConnections"`
		AllowedUpdates []string `yaml:"AllowedUpdates"`
	} `yaml:"TelegramWebHook"`
	YearProgressChatId        int64    `yaml:"YearProgressChatId"`
	DeleteMessageAfterSeconds int64    `yaml:"DeleteMessageAfterSeconds"`
	Currency                  []string `yaml:"Currency"`
}

const config string = `BotToken: "" # telegramBot token
WebhookIP: "127.0.0.1" # webhook ip address
WebhookPort: 6705 # webhook port
TelegramWebHook:    # telegram webhook setting
    Url: ""
    IPAddress: ""
    MaxConnections: 40
    AllowedUpdates: 
        - message
        - edited_message
        - edited_channel_post
        - inline_query
        - chosen_inline_result
        - callback_query
        - shipping_query
        - pre_checkout_query
        - poll
        - poll_answer
        - my_chat_member
        - chat_member
        - chat_join_request
YearProgressChatId: -1
DeleteMessageAfterSeconds: 30
Currency:
    - USD
    - CNY
    - CAD
`

func ReadYaml() (Config, error) {
	filePath := "./config.yml"
	file, _ := os.Stat(filePath)

	var yamlFile []byte
	var err error

	if file == nil {
		fh, _ := os.Create(filePath)
		_, err := fh.WriteString(config)
		if err != nil {
			log.Error(err)
		}
		log.Info("正在生成配置文件, 程序将在5s后退出")
		fh.Close()
		time.Sleep(time.Second * 5)
		os.Exit(0)
	} else {
		yamlFile, err = os.ReadFile(filePath)

	}

	var config Config
	yaml.Unmarshal(yamlFile, &config)
	return config, err
}
