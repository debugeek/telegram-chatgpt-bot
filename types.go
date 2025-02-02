package main

const (
	errTelegramBotTokenNotFound = "telegram bot token not found"

	errFirebaseCredentialNotFound = "firebase credential not found"
	errFirebaseDatabaseNotFound   = "firebase database not found"
)

const (
	CmdSetServiceType    = "setservicetype"
	CmdSetChatGPTAPIKey  = "setchatgptapikey"
	CmdSetChatGPTModel   = "setchatgptmodel"
	CmdSetOllamaEndpoint = "setollamaendpoint"
	CmdSetOllamaModel    = "setollamamodel"
)

const (
	ServiceTypeChatGPT = "chatgpt"
	ServiceTypeOllama  = "ollama"
)

type UserData struct {
	ServiceType    string `firestore:"service-type"`
	ChatGPTModel   string `firestore:"chatgpt-model"`
	ChatGPTAPIKey  string `firestore:"chatgpt-api-key"`
	OllamaEndpoint string `firestore:"ollama-endpoint"`
	OllamaModel    string `firestore:"ollama-model"`
}
