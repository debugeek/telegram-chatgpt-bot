package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Context struct {
	id      int64
	account *Account
	client  *HTTPClient
}

func InitContext() {
	accounts, err := db.GetAccounts()
	if err != nil {
		log.Println(err)
		return
	}

	for _, account := range accounts {
		context := NewContext(account)
		contexts[account.Id] = context
	}
}

func GetCachedContext(id int64, kind int) *Context {
	return contexts[id]
}

func CacheContext(ctx *Context) {
	contexts[ctx.id] = ctx
}

func NewContext(acc *Account) *Context {
	return &Context{
		id:      acc.Id,
		account: acc,
		client: &HTTPClient{
			baseURL: "https://api.openai.com/v1",
		},
	}
}

func (ctx *Context) BuildResponse(msg *tgbotapi.Message) []ChatCompletionMessage {
	var messages []ChatCompletionMessage
	messages = append(messages, ChatCompletionMessage{Role: "user", Content: msg.Text})

	if msg.ReplyToMessage != nil {
		replyToMessageId := msg.ReplyToMessage.MessageID
		for replyToMessageId > 0 {
			replyToMessage, err := db.GetMessage(ctx.account, replyToMessageId)
			if replyToMessage == nil || err != nil {
				break
			}
			messages = append([]ChatCompletionMessage{{Role: replyToMessage.Role, Content: replyToMessage.Text}}, messages...)
			replyToMessageId = replyToMessage.ParentId
		}
	}

	return messages
}

func (ctx *Context) HandleMessage(msg *tgbotapi.Message) string {
	messages := ctx.BuildResponse(msg)

	resp, err := ctx.client.SendChatMessage(messages, ctx.account.Key, ctx.account.Model)
	if err != nil {
		return err.Error()
	}

	return resp.Choices[0].Message.Content
}

func (ctx *Context) SaveMessage(msg *tgbotapi.Message, role string) {
	if msg == nil {
		return
	}

	parentId := -1
	if msg.ReplyToMessage != nil {
		parentId = msg.ReplyToMessage.MessageID
	}
	db.SaveMessage(ctx.account, &Message{
		Id:        msg.MessageID,
		ParentId:  parentId,
		Text:      msg.Text,
		Role:      role,
		Timestamp: msg.Time().Unix(),
	})
}
