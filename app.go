package main

import (
	"encoding/base64"
	"errors"
	"os"
	"telegram-chatgpt-bot/chatgpt"
	"telegram-chatgpt-bot/ollama"

	"github.com/alexflint/go-arg"
	tgbot "github.com/debugeek/telegram-bot"
)

var args struct {
	TelegramBotToken      string `arg:"-t,--tgbot-token" help:"telegram bot token"`
	TelegramBotTokenKey   string `arg:"--tgbot-token-key" help:"env key for telegram bot token"`
	FirebaseCredential    string `arg:"-c,--firebase-credential" help:"base64 encoded firebase credential"`
	FirebaseCredentialKey string `arg:"--firebase-credential-key" help:"env key for base64 encoded firebase credential"`
	FirebaseDatabaseURL   string `arg:"-d,--firebase-database" help:"firebase database url"`
}

type App struct {
	bot      *tgbot.TgBot[UserData]
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

	bot := tgbot.NewBot[UserData](tgbot.Config{
		TelegramBotToken:    telegramBotToken,
		FirebaseCredential:  firebaseCredential,
		FirebaseDatabaseURL: firebaseDatabaseURL,
	})
	bot.RegisterTextHandler(app.processText)
	bot.RegisterCustomCommandHandler(CmdSetServiceType, app.processSetServiceTypeCommand)
	bot.RegisterCustomCommandHandler(CmdSetChatGPTAPIKey, app.processSetChatGPTAPIKeyCommand)
	bot.RegisterCustomCommandHandler(CmdSetChatGPTModel, app.processSetChatGPTModelCommand)
	bot.RegisterCustomCommandHandler(CmdSetOllamaEndpoint, app.processSetOllamaEndpoint)
	bot.RegisterCustomCommandHandler(CmdSetOllamaModel, app.processSetOllamaModelCommand)

	app.bot = bot

	app.firebase = Firebase{
		Firebase: bot.Client.Firebase,
	}

	bot.Start()
}

func (app *App) processText(session tgbot.Session[UserData], text string) {
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

		session.SendText(ollama.Chat(session.User.UserData.OllamaEndpoint, session.User.UserData.OllamaModel, text, 0.5, 200))
	}

}

func (app *App) processSetServiceTypeCommand(session tgbot.Session[UserData], args string) bool {
	session.User.UserData.ServiceType = args
	app.firebase.Firebase.UpdateUser(session.User)
	session.SendText("Service Type is updated.")
	return true
}

func (app *App) processSetChatGPTAPIKeyCommand(session tgbot.Session[UserData], args string) bool {
	session.User.UserData.ChatGPTAPIKey = args
	app.firebase.Firebase.UpdateUser(session.User)
	session.SendText("ChatGPT API Key is updated.")
	return true
}

func (app *App) processSetChatGPTModelCommand(session tgbot.Session[UserData], args string) bool {
	session.User.UserData.ChatGPTModel = args
	app.firebase.Firebase.UpdateUser(session.User)
	session.SendText("ChatGPT Model is updated.")
	return true
}

func (app *App) processSetOllamaEndpoint(session tgbot.Session[UserData], args string) bool {
	session.User.UserData.OllamaEndpoint = args
	app.firebase.Firebase.UpdateUser(session.User)
	session.SendText("Ollama Endpoint is updated.")
	return true
}

func (app *App) processSetOllamaModelCommand(session tgbot.Session[UserData], args string) bool {
	session.User.UserData.OllamaModel = args
	app.firebase.Firebase.UpdateUser(session.User)
	session.SendText("Ollama Model is updated.")
	return true
}
