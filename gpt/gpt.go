package gpt

import (
	"fmt"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jellyqwq/Paimon/chatgpt"
	pconfig "github.com/jellyqwq/Paimon/config"

	"github.com/m1guelpf/chatgpt-telegram/src/markdown"
	"github.com/m1guelpf/chatgpt-telegram/src/ratelimit"
	"github.com/m1guelpf/chatgpt-telegram/src/session"
)

type Conversation struct {
	ConversationID string
	LastMessageID  string
}
type GPT struct {
	mutex             sync.RWMutex
	userConversations map[int64]*Conversation
	gpt               chatgpt.ChatGPT
}

func New(configGPT pconfig.GPTConfig) *GPT {
	if configGPT.OpenAISession == "" {
		session, err := session.GetSession()
		if err != nil {
			log.Printf("Couldn't get OpenAI session: %v", err)
		}

		err = configGPT.Set("OpenAISession", session)
		if err != nil {
			log.Printf("Couldn't save OpenAI session: %v", err)
		}
	}
	g, err := chatgpt.Init(configGPT)
	if err != nil {
		log.Fatal(err)
	}
	return &GPT{
		userConversations: map[int64]*Conversation{},
		gpt:               g,
	}
}

func (g *GPT) messageExists(id int64) (*Conversation, bool) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	e, ok := g.userConversations[id]
	return e, ok
}

func (g *GPT) messageAdd(id int64) *Conversation {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	c := &Conversation{}
	g.userConversations[id] = c
	return c
}

func (g *GPT) conversationModify(id int64, cid, lastid string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.userConversations[id].ConversationID = cid
	g.userConversations[id].LastMessageID = lastid
}

func (g *GPT) NewMessage(bot *tgbotapi.BotAPI, Message *tgbotapi.Message) {
	entry, ok := g.messageExists(Message.Chat.ID)
	if !ok {
		entry = g.messageAdd(Message.Chat.ID)
	}
	msg := tgbotapi.NewMessage(Message.Chat.ID, "")
	msg.ReplyToMessageID = Message.MessageID
	msg.ParseMode = "Markdown"
	bot.Request(tgbotapi.NewChatAction(Message.Chat.ID, "typing"))
	feed, err := chatGPT.SendMessage(Message.Text, entry.ConversationID, entry.LastMessageID)
	if err != nil {
		msg.Text = fmt.Sprintf("Error: %v", err)
	}
	var message tgbotapi.Message
	var lastResp string

	debouncedType := ratelimit.Debounce((10 * time.Second), func() {
		bot.Request(tgbotapi.NewChatAction(Message.Chat.ID, "typing"))
	})
	debouncedEdit := ratelimit.DebounceWithArgs((1 * time.Second), func(text interface{}, messageId interface{}) {
		_, err = bot.Request(tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    msg.ChatID,
				MessageID: messageId.(int),
			},
			Text:      text.(string),
			ParseMode: "Markdown",
		})

		if err != nil {
			if err.Error() == "Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message" {
				return
			}

			log.Printf("Couldn't edit message: %v", err)
		}
	})
pollResponse:
	for {
		debouncedType()

		response, ok := <-feed
		if !ok {
			break pollResponse
		}

		g.conversationModify(Message.Chat.ID, response.MessageId, response.ConversationId)

		lastResp = markdown.EnsureFormatting(response.Message)
		msg.Text = lastResp

		if message.MessageID == 0 {
			message, err = bot.Send(msg)
			if err != nil {
				log.Printf("Couldn't send message: %v", err)
			}
		} else {
			debouncedEdit(lastResp, message.MessageID)
		}

		_, err = bot.Request(tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    msg.ChatID,
				MessageID: message.MessageID,
			},
			Text:      lastResp,
			ParseMode: "Markdown",
		})

		if err != nil {
			if err.Error() == "Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message" {
				continue
			}

			log.Printf("Couldn't perform final edit on message: %v", err)
		}

		continue
	}
}
