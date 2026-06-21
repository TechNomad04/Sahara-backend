package app

import (
	"sahara/internal/store"

	"gorm.io/gorm"
)

type Dependencies struct {
	DB    *gorm.DB
	Redis *store.Redis
	Qdrant *store.Qdrant
}