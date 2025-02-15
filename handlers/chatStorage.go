package handlers

import "sync"

var (
	ChatMutex   sync.Mutex
	ActiveChats = make(map[string]string) // user_id -> chat_id
)
