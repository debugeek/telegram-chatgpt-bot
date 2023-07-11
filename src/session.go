package main

import (
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Session struct {
	bot   *tgbotapi.BotAPI
	token string
}

func InitSession() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	_, err = bot.GetUpdates(u)
	if err != nil {
		log.Fatal(err)
	}

	session = &Session{
		token: token,
		bot:   bot,
	}
	session.Run()
	log.Println(`Session initialized`)
}

func (s *Session) Run() {
	go s.RunLoop()
}

func (s *Session) RunLoop() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates, err := s.bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
		return
	}

	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		var message *tgbotapi.Message
		if update.Message != nil {
			message = update.Message
		} else if update.ChannelPost != nil {
			message = update.ChannelPost
		}
		if message == nil {
			log.Println("unabled to handle update")
			continue
		}

		id := message.Chat.ID
		kind := -1
		if message.Chat.IsPrivate() {
			kind = 0
		} else if message.Chat.IsGroup() || message.Chat.IsSuperGroup() {
			kind = 1
		} else if message.Chat.IsChannel() {
			kind = 2
		}

		context := GetCachedContext(id, kind)
		if context == nil {
			account := &Account{
				Id:     id,
				Kind:   kind,
				Status: 1,
			}
			err = db.SaveAccount(account)
			if err != nil {
				log.Println(err)
				continue
			}
			context = NewContext(account)
			CacheContext(context)
		}

		s.handleMessage(context, message)
	}
}

func (s *Session) Send(ctx *Context, text string, disableWebPagePreview bool) (*tgbotapi.Message, error) {
	if ctx.account.Status == -1 {
		return nil, nil
	}

	cfg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           ctx.id,
			ReplyToMessageID: 0,
		},
		Text:                  text,
		ParseMode:             "markdown",
		DisableWebPagePreview: disableWebPagePreview,
	}
	msg, err := s.bot.Send(cfg)
	if err != nil {
		s.handleError(ctx, err)
	}
	return &msg, err
}

func (s *Session) Reply(ctx *Context, replyToMessageID int, text string, disableWebPagePreview bool) (*tgbotapi.Message, error) {
	if ctx.account.Status == -1 {
		return nil, nil
	}

	cfg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           ctx.id,
			ReplyToMessageID: replyToMessageID,
		},
		Text:                  text,
		ParseMode:             "markdown",
		DisableWebPagePreview: disableWebPagePreview,
	}
	msg, err := s.bot.Send(cfg)
	if err != nil {
		s.handleError(ctx, err)
	}
	return &msg, err
}

func (s *Session) handleMessage(ctx *Context, msg *tgbotapi.Message) {
	if ctx.account.Status == -1 {
		ctx.account.Status = 1
		db.SaveAccount(ctx.account)
	}

	replyToMessageID := msg.MessageID

	if msg.IsCommand() {
		switch strings.ToLower(msg.Command()) {
		case "start":
			{
				s.Reply(ctx, replyToMessageID, "Greetings.", false)
				break
			}
		case "setapikey":
			{
				args := msg.CommandArguments()
				ctx.account.APIKey = args
				db.SaveAccount(ctx.account)
				s.Reply(ctx, replyToMessageID, "API Key is updated.", false)
			}
		case "setmodel":
			{
				args := msg.CommandArguments()
				ctx.account.Model = args
				db.SaveAccount(ctx.account)
				s.Reply(ctx, replyToMessageID, "Model is updated.", false)
			}
		default:
			break
		}
	} else {
		response := ctx.HandleMessage(msg)
		s.Reply(ctx, replyToMessageID, response, false)
	}
}

func (s *Session) handleError(ctx *Context, err error) {
	switch err.Error() {
	case errChatNotFound, errNotMember:
		ctx.account.Status = -1
		db.SaveAccount(ctx.account)
	}
}
