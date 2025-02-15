package models

import (
	"gorm.io/gorm"
)

// Структура сообщений
type Message struct {
	gorm.Model
	ChatID    string `json:"chat_id"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}
