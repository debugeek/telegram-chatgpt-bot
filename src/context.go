package main

import "log"

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

func CacheContext(context *Context) {
	contexts[context.id] = context
}

func NewContext(account *Account) *Context {
	return &Context{
		id:      account.Id,
		account: account,
		client: &HTTPClient{
			baseURL: "https://api.openai.com/v1",
		},
	}
}

// Message Handler

func (ctx *Context) HandleMessage(message string) string {
	if len(ctx.account.Key) == 0 {
		return "API key missing"
	}

	if len(ctx.account.Model) == 0 {
		return "Model missing"
	}

	resp, err := ctx.client.SendChatMessage(message, ctx.account.Key, ctx.account.Model)
	if err != nil {
		return err.Error()
	}

	return resp.Choices[0].Message.Content
}
