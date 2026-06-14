package models

import (
	"time"

	"github.com/lib/pq"
)

type RequestStatus string

const (
	StatusPending RequestStatus = "OPEN"
	StatusClosed  RequestStatus = "CLOSED"
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
	ID               uint           `gorm:"primaryKey" json:"id"`
	Title            string         `gorm:"not null" json:"title"`
	Description      string         `gorm:"not null" json:"description"`
	Categories       pq.StringArray `gorm:"type:text[]" json:"categories"`
	Country          string         `gorm:"size:100;not null" json:"country"`
	State            string         `gorm:"size:100;not null" json:"state"`
	City             string         `gorm:"size:100;not null" json:"city"`
	Status           RequestStatus  `gorm:"type:varchar(20);default:'working'" json:"status"`
	IssuerId uint `json:"issuerId"`
	Issuer Issuer `gorm:"foreignKey:IssuerId"`
	ParticipantCount int            `gorm:"not null" json:"participantCount"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"createdAt"`
}
