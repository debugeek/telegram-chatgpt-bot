package main

const (
	errTelegramBotTokenNotFound = "telegram bot token not found"

	errFirebaseCredentialNotFound = "firebase credential not found"
	errFirebaseDatabaseNotFound   = "firebase database not found"
)

const (
	CmdSetServiceType    = "setservicetype"
	CmdSetup             = "setup"
	CmdSetChatGPTAPIKey  = "setchatgptapikey"
	CmdSetChatGPTModel   = "setchatgptmodel"
	CmdSetOllamaEndpoint = "setollamaendpoint"
	CmdSetOllamaAPIKey   = "setollamaapikey"
	CmdSetOllamaModel    = "setollamamodel"
)

const (
	ServiceTypeChatGPT = "chatgpt"
	ServiceTypeOllama  = "ollama"
)

type BotData struct {
}

type UserData struct {
	ServiceType    string `firestore:"service-type"`
	ChatGPTModel   string `firestore:"chatgpt-model"`
	ChatGPTAPIKey  string `firestore:"chatgpt-api-key"`
	OllamaEndpoint string `firestore:"ollama-endpoint"`
	OllamaAPIKey   string `firestore:"ollama-api-key"`
	OllamaModel    string `firestore:"ollama-model"`
}
