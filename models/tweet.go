package models

import "gorm.io/gorm"

type Tweet struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	UserID  uint   `gorm:"not null"`
	User    User   `gorm:"constraint:OnDelete:CASCADE;" json:"user"`
}
