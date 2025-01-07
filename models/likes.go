package models

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	TweetID uint  `gorm:"not null"`
	UserID  uint  `gorm:"not null"`
	Tweet   Tweet `gorm:"constraint:OnDelete:CASCADE;"`
	User    User  `gorm:"constraint:OnDelete:CASCADE;"`
}
