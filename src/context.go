package main

import (
	"log"

	chatgpt "telegram-chatgpt-bot/src/chatgpt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Context struct {
	id      int64
	account *Account
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
	}
}

func (ctx *Context) HandleMessage(msg *tgbotapi.Message) string {
	if ctx.account.APIKey == "" {
		return "API Key is missing."
	} else if ctx.account.Model == "" {
		return "Model is missing."
	} else {
		return chatgpt.SendText(msg.Text, ctx.account.APIKey, ctx.account.Model)
	}
}
