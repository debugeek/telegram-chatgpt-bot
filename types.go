package main

const (
	errTelegramBotTokenNotFound = "telegram bot token not found"

	errFirebaseCredentialNotFound = "firebase credential not found"
	errFirebaseDatabaseNotFound   = "firebase database not found"
)

const (
	CmdSetAPIKey = "setapikey"
	CmdSetModel  = "setmodel"
)

type UserData struct {
	Model  string `firestore:"chatgpt-model"`
	APIKey string `firestore:"chatgpt-api-key"`
}
