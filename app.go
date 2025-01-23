package main

import (
	"encoding/base64"
	"errors"
	"os"
	"telegram-chatgpt-bot/chatgpt"

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
	bot.RegisterCustomCommandHandler(CmdSetAPIKey, app.processSetAPIKeyCommand)
	bot.RegisterCustomCommandHandler(CmdSetModel, app.processSetModelCommand)

	app.bot = bot

	app.firebase = Firebase{
		Firebase: bot.Client.Firebase,
	}

	bot.Start()
}

func (app *App) processText(session tgbot.Session[UserData], text string) {
	if session.User.UserData.APIKey == "" {
		session.SendText("API Key is missing.")
		return
	}
	if session.User.UserData.Model == "" {
		session.SendText("Model is missing.")
		return
	}

	session.SendText(chatgpt.Chat(text, session.User.UserData.APIKey, session.User.UserData.Model))
}

func (app *App) processSetAPIKeyCommand(session tgbot.Session[UserData], args string) bool {
	session.User.UserData.APIKey = args
	app.firebase.Firebase.UpdateUser(session.User)
	session.SendText("API Key is updated.")
	return true
}

func (app *App) processSetModelCommand(session tgbot.Session[UserData], args string) bool {
	session.User.UserData.Model = args
	app.firebase.Firebase.UpdateUser(session.User)
	session.SendText("Model is updated.")
	return true
}
