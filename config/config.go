package config

import (
	"log"
	"os"
    "time"

	"gopkg.in/yaml.v2"
)

type Config struct {
    BotToken	string	`yaml:"BotToken"`
    WebhookIP    string  `yaml:"WebhookIP"`
    WebhookPort     uint64  `yaml:"WebhookPort"`
    TelegramWebHook		struct 	{
        Url			string `yaml:"Url"`
        IPAddress	string `yaml:"IPAddress"`
        MaxConnections int `yaml:"MaxConnections"`
        AllowedUpdates []string `yaml:"AllowedUpdates"`
    } `yaml:"TelegramWebHook"`
    CQ2TG   struct {
        MentionReflection map[uint64]uint64 `yaml:"MentionReflection"`
        RecivedChatId int64 `yaml:"RecivedChatId"`
    } `yaml:"CQ2TG"`
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