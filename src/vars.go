package main

var (
	token    string
	session  *Session
	db       *Firebase
	contexts map[int64]*Context = make(map[int64]*Context)
)

const (
	errChatNotFound string = "Bad Request: chat not found"
	errNotMember    string = "Forbidden: bot is not a member of the channel chat"
)
