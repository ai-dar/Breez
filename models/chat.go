package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	UserID string `json:"user_id"`
	ChatID string `json:"chat_id"`
	Active bool   `json:"active" gorm:"default:true"`
}
