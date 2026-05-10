package main

import (
	"strings"

	tgbot "github.com/debugeek/telegram-bot"
)

func (app *App) processSetupCommand(session *tgbot.Session[BotData, UserData], args string, message *tgbot.Message) tgbot.CmdResult {
	args = strings.TrimSpace(args)
	if strings.EqualFold(args, setupButtonCancel) {
		app.sendTextWithRemoveKeyboard(session, "Setup is cancelled.", message.MessageID)
		return tgbot.CmdResultProcessed
	}

	switch session.CommandSession.Stage {
	case "":
		session.CommandSession.Stage = setupStageService
		app.sendTextWithKeyboard(session, "Choose a service to configure.", message.MessageID, setupServiceKeyboard())
		return tgbot.CmdResultWaitingForInput
	case setupStageService:
		return app.processSetupService(session, args, message)
	default:
		return app.processSetupInput(session, args, message)
	}
}

func (app *App) processSetupService(session *tgbot.Session[BotData, UserData], args string, message *tgbot.Message) tgbot.CmdResult {
	switch args {
	case setupButtonChatGPT:
		session.User.UserData.ServiceType = ServiceTypeChatGPT
		session.CommandSession.Stage = setupStageChatGPTAPIKey
		app.firebase.Firebase.UpdateUser(session.User)
		app.sendTextWithKeyboard(session, "Send your ChatGPT API key.", message.MessageID, setupCancelKeyboard())
		return tgbot.CmdResultWaitingForInput
	case setupButtonOllamaCloud:
		session.User.UserData.ServiceType = ServiceTypeOllama
		session.User.UserData.OllamaEndpoint = "https://ollama.com/api/chat"
		session.CommandSession.Stage = setupStageOllamaAPIKey
		app.firebase.Firebase.UpdateUser(session.User)
		app.sendTextWithKeyboard(session, "Send your Ollama API key.", message.MessageID, setupCancelKeyboard())
		return tgbot.CmdResultWaitingForInput
	case setupButtonOllamaLocal:
		session.User.UserData.ServiceType = ServiceTypeOllama
		session.User.UserData.OllamaAPIKey = ""
		session.CommandSession.Stage = setupStageOllamaEndpoint
		app.firebase.Firebase.UpdateUser(session.User)
		app.sendTextWithKeyboard(session, "Choose or send your local Ollama endpoint.", message.MessageID, setupOllamaLocalEndpointKeyboard())
		return tgbot.CmdResultWaitingForInput
	default:
		app.sendTextWithKeyboard(session, "Please choose one of the setup options.", message.MessageID, setupServiceKeyboard())
		return tgbot.CmdResultWaitingForInput
	}
}

func (app *App) processSetupInput(session *tgbot.Session[BotData, UserData], args string, message *tgbot.Message) tgbot.CmdResult {
	step, ok := setupInputSteps[session.CommandSession.Stage]
	if !ok {
		session.CommandSession.Stage = setupStageService
		app.sendTextWithKeyboard(session, "Choose a service to configure.", message.MessageID, setupServiceKeyboard())
		return tgbot.CmdResultWaitingForInput
	}

	if args == "" {
		app.sendTextWithKeyboard(session, step.Prompt, message.MessageID, setupCancelKeyboard())
		return tgbot.CmdResultWaitingForInput
	}

	step.Apply(&session.User.UserData, args)
	app.firebase.Firebase.UpdateUser(session.User)
	if step.NextStage == "" {
		app.sendTextWithRemoveKeyboard(session, step.NextPrompt, message.MessageID)
		return tgbot.CmdResultProcessed
	}

	session.CommandSession.Stage = step.NextStage
	app.sendTextWithKeyboard(session, step.NextPrompt, message.MessageID, setupCancelKeyboard())
	return tgbot.CmdResultWaitingForInput
}

const (
	setupStageService        = "service"
	setupStageChatGPTAPIKey  = "chatgpt-api-key"
	setupStageChatGPTModel   = "chatgpt-model"
	setupStageOllamaAPIKey   = "ollama-api-key"
	setupStageOllamaEndpoint = "ollama-endpoint"
	setupStageOllamaModel    = "ollama-model"
	setupButtonChatGPT       = "ChatGPT"
	setupButtonOllamaCloud   = "Ollama Cloud"
	setupButtonOllamaLocal   = "Ollama Local"
	setupButtonCancel        = "Cancel Setup"
)

type setupInputStep struct {
	Prompt     string
	NextStage  string
	NextPrompt string
	Apply      func(*UserData, string)
}

var setupInputSteps = map[string]setupInputStep{
	setupStageChatGPTAPIKey: {
		Prompt:     "Send your ChatGPT API key.",
		NextStage:  setupStageChatGPTModel,
		NextPrompt: "Send the ChatGPT model name.",
		Apply: func(userData *UserData, value string) {
			userData.ChatGPTAPIKey = value
		},
	},
	setupStageChatGPTModel: {
		Prompt:     "Send the ChatGPT model name.",
		NextPrompt: "ChatGPT setup is complete.",
		Apply: func(userData *UserData, value string) {
			userData.ChatGPTModel = value
		},
	},
	setupStageOllamaAPIKey: {
		Prompt:     "Send your Ollama API key.",
		NextStage:  setupStageOllamaModel,
		NextPrompt: "Send the Ollama model name.",
		Apply: func(userData *UserData, value string) {
			userData.OllamaAPIKey = value
		},
	},
	setupStageOllamaModel: {
		Prompt:     "Send the Ollama model name.",
		NextPrompt: "Ollama setup is complete.",
		Apply: func(userData *UserData, value string) {
			userData.OllamaModel = value
		},
	},
	setupStageOllamaEndpoint: {
		Prompt:     "Send your local Ollama endpoint.",
		NextStage:  setupStageOllamaModel,
		NextPrompt: "Send the Ollama model name.",
		Apply: func(userData *UserData, value string) {
			userData.OllamaEndpoint = value
		},
	},
}

func (app *App) sendTextWithKeyboard(session *tgbot.Session[BotData, UserData], text string, replyToMessageID int, keyboard tgbot.ReplyKeyboardMarkup) error {
	return session.SendTextWithConfig(text, tgbot.MessageConfig{
		ReplyToMessageID: replyToMessageID,
		ReplyMarkup:      keyboard,
	})
}

func (app *App) sendTextWithRemoveKeyboard(session *tgbot.Session[BotData, UserData], text string, replyToMessageID int) error {
	return session.SendTextWithConfig(text, tgbot.MessageConfig{
		ReplyToMessageID: replyToMessageID,
		ReplyMarkup:      tgbot.NewRemoveKeyboard(false),
	})
}

func setupServiceKeyboard() tgbot.ReplyKeyboardMarkup {
	keyboard := tgbot.NewReplyKeyboard(
		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton(setupButtonChatGPT),
			tgbot.NewKeyboardButton(setupButtonOllamaCloud),
		),
		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton(setupButtonOllamaLocal),
			tgbot.NewKeyboardButton(setupButtonCancel),
		),
	)
	keyboard.OneTimeKeyboard = true
	return keyboard
}

func setupCancelKeyboard() tgbot.ReplyKeyboardMarkup {
	keyboard := tgbot.NewReplyKeyboard(
		tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton(setupButtonCancel)),
	)
	keyboard.OneTimeKeyboard = true
	return keyboard
}

func setupOllamaLocalEndpointKeyboard() tgbot.ReplyKeyboardMarkup {
	keyboard := tgbot.NewReplyKeyboard(
		tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton("http://localhost:11434/api/chat")),
		tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton(setupButtonCancel)),
	)
	keyboard.OneTimeKeyboard = true
	return keyboard
}
