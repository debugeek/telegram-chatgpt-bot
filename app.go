package main

import (
	"encoding/base64"
	"errors"
	"os"
	"telegram-chatgpt-bot/chatgpt"
	"telegram-chatgpt-bot/ollama"

	"github.com/alexflint/go-arg"
	tgbot "github.com/debugeek/telegram-bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var args struct {
	TelegramBotToken      string `arg:"-t,--tgbot-token" help:"telegram bot token"`
	TelegramBotTokenKey   string `arg:"--tgbot-token-key" help:"env key for telegram bot token"`
	FirebaseCredential    string `arg:"-c,--firebase-credential" help:"base64 encoded firebase credential"`
	FirebaseCredentialKey string `arg:"--firebase-credential-key" help:"env key for base64 encoded firebase credential"`
	FirebaseDatabaseURL   string `arg:"-d,--firebase-database" help:"firebase database url"`
}

type App struct {
	bot      *tgbot.TgBot[BotData, UserData]
	firebase Firebase
}

func (app *App) launch() {
	arg.MustParse(&args)

	telegramBotToken := args.TelegramBotToken
	if telegramBotToken == "" {
		telegramBotToken = os.Getenv(args.TelegramBotTokenKey)
	}
	if telegramBotToken == "" {
		panic(errors.New(errTelegramBotTokenNotFound))
	}

	encodedFirebaseCredential := args.FirebaseCredential
	if encodedFirebaseCredential == "" {
		encodedFirebaseCredential = os.Getenv(args.FirebaseCredentialKey)
	}
	if encodedFirebaseCredential == "" {
		panic(errors.New(errFirebaseCredentialNotFound))
	}
	firebaseCredential, err := base64.StdEncoding.DecodeString(encodedFirebaseCredential)
	if err != nil {
		panic(err)
	}

	firebaseDatabaseURL := args.FirebaseDatabaseURL
	if firebaseDatabaseURL == "" {
		panic(errors.New(errFirebaseDatabaseNotFound))
	}

	bot := tgbot.NewBot[BotData, UserData](tgbot.Config{
		TelegramBotToken:    telegramBotToken,
		FirebaseCredential:  firebaseCredential,
		FirebaseDatabaseURL: firebaseDatabaseURL,
	}, app)
	bot.RegisterTextHandler(app.processText)
	bot.RegisterCommandHandler(CmdSetServiceType, app.processSetServiceTypeCommand)
	bot.RegisterCommandHandler(CmdSetChatGPTAPIKey, app.processSetChatGPTAPIKeyCommand)
	bot.RegisterCommandHandler(CmdSetChatGPTModel, app.processSetChatGPTModelCommand)
	bot.RegisterCommandHandler(CmdSetOllamaEndpoint, app.processSetOllamaEndpointCommand)
	bot.RegisterCommandHandler(CmdSetOllamaModel, app.processSetOllamaModelCommand)

	app.bot = bot

	app.firebase = Firebase{
		Firebase: bot.Client.Firebase,
	}

	bot.Start()
}

func (app *App) NewUserData() UserData {
	return UserData{}
}

func (app *App) DidLoadUser(session *tgbot.Session[BotData, UserData], user *tgbot.User[UserData]) {

}

func (app *App) DidLoadPreference() {

}

func (app *App) processText(session *tgbot.Session[BotData, UserData], text string, message *tgbotapi.Message) {
	switch session.User.UserData.ServiceType {
	case "", ServiceTypeChatGPT:
		if session.User.UserData.ChatGPTAPIKey == "" {
			session.SendText("ChatGPT API Key is missing.")
			return
		}
		if session.User.UserData.ChatGPTModel == "" {
			session.SendText("ChatGPT Model is missing.")
			return
		}

		session.SendText(chatgpt.Chat(text, session.User.UserData.ChatGPTAPIKey, session.User.UserData.ChatGPTModel))

	case ServiceTypeOllama:
		if session.User.UserData.OllamaEndpoint == "" {
			session.SendText("Ollama Endpoint is missing.")
			return
		}
		if session.User.UserData.OllamaModel == "" {
			session.SendText("Ollama Model is missing.")
			return
		}

		session.SendTextWithConfig(ollama.Chat(session.User.UserData.OllamaEndpoint, session.User.UserData.OllamaModel, text, 0.5, 200), tgbot.MessageConfig{
			ReplyToMessageID: message.MessageID,
			ParseMode:        tgbot.ParseModeMarkdown,
		})
	}

}

func (app *App) processSetServiceTypeCommand(session *tgbot.Session[BotData, UserData], args string, message *tgbotapi.Message) tgbot.CmdResult {
	session.User.UserData.ServiceType = args
	app.firebase.Firebase.UpdateUser(session.User)
	session.SendTextWithConfig("Service Type is updated.", tgbot.MessageConfig{
		ReplyToMessageID: message.MessageID,
	})
	return tgbot.CmdResultProcessed
}

func (app *App) processSetChatGPTAPIKeyCommand(session *tgbot.Session[BotData, UserData], args string, message *tgbotapi.Message) tgbot.CmdResult {
	session.User.UserData.ChatGPTAPIKey = args
	app.firebase.Firebase.UpdateUser(session.User)
	session.SendTextWithConfig("ChatGPT API Key is updated.", tgbot.MessageConfig{
		ReplyToMessageID: message.MessageID,
	})
	return tgbot.CmdResultProcessed
}

func (app *App) processSetChatGPTModelCommand(session *tgbot.Session[BotData, UserData], args string, message *tgbotapi.Message) tgbot.CmdResult {
	session.User.UserData.ChatGPTModel = args
	app.firebase.Firebase.UpdateUser(session.User)
	session.SendTextWithConfig("ChatGPT Model is updated.", tgbot.MessageConfig{
		ReplyToMessageID: message.MessageID,
	})
	return tgbot.CmdResultProcessed
}

func (app *App) processSetOllamaEndpointCommand(session *tgbot.Session[BotData, UserData], args string, message *tgbotapi.Message) tgbot.CmdResult {
	session.User.UserData.OllamaEndpoint = args
	app.firebase.Firebase.UpdateUser(session.User)
	session.SendTextWithConfig("Ollama Endpoint is updated.", tgbot.MessageConfig{
		ReplyToMessageID: message.MessageID,
	})
	return tgbot.CmdResultProcessed
}

func (app *App) processSetOllamaModelCommand(session *tgbot.Session[BotData, UserData], args string, message *tgbotapi.Message) tgbot.CmdResult {
	session.User.UserData.OllamaModel = args
	app.firebase.Firebase.UpdateUser(session.User)
	session.SendTextWithConfig("Ollama Model is updated.", tgbot.MessageConfig{
		ReplyToMessageID: message.MessageID,
	})
	return tgbot.CmdResultProcessed
}
