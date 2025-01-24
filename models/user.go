package models

type User struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Email      string `gorm:"unique" json:"email"`
	Password   string `json:"password"`
	Name       string `json:"name,omitempty"`
	Role       string `gorm:"default:'user'"`
	IsVerified bool   `gorm:"default:false"`
}
