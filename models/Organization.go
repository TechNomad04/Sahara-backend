package models

type Organization struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"not null" json:"email"`
	GoogleId string `gorm:"not null" json:"googleId"`
}
