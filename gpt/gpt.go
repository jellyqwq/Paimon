package gpt

import (
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jellyqwq/Paimon/chatgpt"
	pconfig "github.com/jellyqwq/Paimon/config"

	"github.com/m1guelpf/chatgpt-telegram/src/markdown"
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
	feed, err := g.gpt.SendMessage(Message.Text, entry.ConversationID, entry.LastMessageID)
	if err != nil {
		msg.Text = fmt.Sprintf("Error: %v", err)
	}
	//var message tgbotapi.Message
	var lastResp string
	seq := ""
	for {
		response, ok := <-feed
		if !ok {
			g.conversationModify(Message.Chat.ID, response.ConversationId, response.MessageId)
			break
		}
		seq += response.Message
	}

	lastResp = markdown.EnsureFormatting(seq)
	msg.Text = lastResp
	_, err = bot.Send(msg)

	if err != nil {
		log.Printf("Couldn't perform final edit on message: %v", err)
	}

}
