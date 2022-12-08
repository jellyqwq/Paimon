package config

import (
	"log"
	"os"
	"time"
	"fmt"

	"gopkg.in/yaml.v2"
	"github.com/spf13/viper"
)

type GPTConfig struct {
	OpenAISession string
}

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
	CQ2TG struct {
		MentionReflection map[uint64]uint64 `yaml:"MentionReflection"`
		RecivedChatId     int64             `yaml:"RecivedChatId"`
	} `yaml:"CQ2TG"`
	DeleteMessageAfterSeconds int64    `yaml:"DeleteMessageAfterSeconds"`
	Currency                  []string `yaml:"Currency"`
	GPTChatid string `yaml:"GPTChatid"`
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
CQ2TG:
    MentionReflection: 
        # QQ Number: Telegram User ID
        00000000: 11111111
    RecivedChatId: 666 # chat id to recive message
DeleteMessageAfterSeconds: 10
Currency:
    - USD
    - CNY
    - CAD
GPTChatid: ""
`

func ReadYaml() (Config, error) {
	filePath := "./config.yml"
	file, _ := os.Stat(filePath)

	var yamlFile []byte
	var err error

	if file == nil {
		fh, _ := os.Create(filePath)
		n, err := fh.WriteString(config)
		log.Println(n, err)
		log.Println("正在生成配置文件, 程序将在5s后退出")
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

// =============== OpenAi-GPT ==============
// Copy from https://github.com/m1guelpf/chatgpt-telegram/blob/main/src/config/config.go
// init tries to read the config from the file, and creates it if it doesn't exist.
func GPTinit() (GPTConfig, error) {
	viper.SetConfigType("json")
	viper.SetConfigName("chatgpt")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				return GPTConfig{}, fmt.Errorf("couldn't create config file: %v", err)
			}
		} else {
			return GPTConfig{}, fmt.Errorf("couldn't read config file: %v", err)
		}
	}

	var cfg GPTConfig
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return GPTConfig{}, fmt.Errorf("error parsing config: %v", err)
	}

	return cfg, nil
}

// key should be part of the Config struct
func (cfg *GPTConfig) Set(key string, value interface{}) error {
	viper.Set(key, value)

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return fmt.Errorf("error parsing config: %v", err)
	}

	return viper.WriteConfig()
}
// =============== OpenAi-GPT ==============