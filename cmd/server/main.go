package main

import (
	"rio/internal/handlers"
	"rio/internal/setup"
	"rio/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	deps := setup.Setup()

	router := gin.Default()

	public := router.Group("/api")
	// Users
	public.POST("/register", deps.UserHandler.Register)
	public.POST("/login", deps.UserHandler.Login)
	public.GET("/users", deps.UserHandler.GetUsers)
	public.GET("/users/:username", deps.UserHandler.FindUsername)

	//Admin
	protected := router.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", deps.UserHandler.CurrentUser)

	// Server
	router.POST("/servers", handlers.CreateServer)
	router.GET("/servers", handlers.GetServers)

	// Channel
	router.POST("/channels", handlers.CreateChannel)
	router.GET("/channels/:server_id", handlers.GetChannels)

	// Messages
	router.POST("/messages/:channel_id", handlers.SendMessage)
	router.GET("/messages/:channel_id", handlers.GetMessages)

	router.Run("localhost:8080")

}
