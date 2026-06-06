package main

import (
	"sahara/db"
	"sahara/internal/auth"
	"sahara/internal/store"

	"github.com/gin-gonic/gin"
)

func main() {

	db := db.InitDB()
	redisClient := store.NewRedis()

	h := &auth.Handler{
		DB: db,
		Redis: redisClient,
	}

	r := gin.Default()
	r.GET("/auth/google", auth.GoogleLogin)
	r.GET("/auth/google/callback", h.GoogleCallback)
	r.POST("/auth/refresh", h.Refresh)
	r.POST("/auth/logout",
		auth.AuthMiddleware(redisClient),
		h.Logout,
	)

	r.Run(":8080")
}