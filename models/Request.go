package models

import (
	"time"
)

type RequestStatus string

const (
	StatusPending  RequestStatus = "OPEN"
	StatusClosed RequestStatus = "CLOSED"
	StatusWorking RequestStatus = "WORKING"
)

type CategoryTypes string

const (
	Healthcare     CategoryTypes = "healthcare"
	Education      CategoryTypes = "education"
	Environment    CategoryTypes = "environment"
	Community      CategoryTypes = "community"
	DisasterRelief CategoryTypes = "disaster-relief"
	AnimalWelfare  CategoryTypes = "animal-welfare"
	Others         CategoryTypes = "others"
)

type Request struct {
	ID               uint          `gorm:"primary key" json:"id"`
	Title            string        `gorm:"not null" json:"title"`
	Description      string        `gorm:"not null" json:"description"`
	Category         CategoryTypes `gorm:"type:varchar(20);default:'others'" json:"category"`
	Location         string        `gorm:"not null" json:"location"`
	Status           RequestStatus `gorm:"type:varchar(20);default:'working'" json:"status"`
	IssuerName       string        `gorm:"not null" json:"issuerName"`
	IssuerType       string        `gorm:"not null" json:"issuerType"`
	ParticipantCount int           `gorm:"type:INTEGER;not null" json:"participantCount"`
	CreatedAt        time.Time     `gorm:"not null" json:"createdAt"`
}
