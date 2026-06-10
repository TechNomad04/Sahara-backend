package models

type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	GoogleId string `gorm:"unique;not null" json:"googleId"`
	Name  string `json:"name"`
	Email string `gorm:"unique" json:"email"`
}