package main

import "log"

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
		context := &Context{
			id:      account.Id,
			account: account,
		}

		contexts[account.Id] = context
	}
}

func GetCachedContext(id int64, kind int) *Context {
	return contexts[id]
}

func CacheContext(context *Context) {
	contexts[context.id] = context
}

// Message Handler

func (context *Context) HandleMessage(message string) string {
	if len(context.account.Key) == 0 {
		return "API key missing"
	}

	if len(context.account.Model) == 0 {
		return "Model missing"
	}

	return message
}
