package main

import (
	"log"
	"sahara/db"
	"sahara/internal/app"
	"sahara/internal/auth"
	"sahara/internal/store"
	"sahara/models"

	"github.com/gin-gonic/gin"
)

func main() {
	database := db.InitDB()
	redisClient := store.NewRedis()

	deps := &app.Dependencies{
		DB:    database,
		Redis: redisClient,
	}

	h := &auth.Handler{
		Dependencies: deps,
	}

	r := gin.Default()

	err := database.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatal("failed to migrate database:", err)
    }

	r.GET("/auth/google", auth.GoogleLogin)
	r.GET("/auth/google/callback", h.GoogleCallback)

	r.POST("/auth/refresh", h.Refresh)

	r.POST(
		"/auth/logout",
		auth.AuthMiddleware(redisClient),
		h.Logout,
	)

	r.Run(":8080")
}