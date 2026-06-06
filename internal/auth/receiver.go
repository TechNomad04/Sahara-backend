package auth

import (
	"database/sql"
	"sahara/internal/store"
)

type Handler struct {
	DB *sql.DB
	Redis *store.Redis
}