package models

type Issuer struct {
	ID uint `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
	IssuerType string `gorm:"not null" jsoN:"issuerType"`
}