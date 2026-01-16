package main

import (
	"rio/internal/setup"
	"rio/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	deps := setup.Setup()

	router := gin.Default()

	public := router.Group("/api")

	public.POST("/register", deps.UserHandler.Register)
	public.POST("/login", deps.UserHandler.Login)
	public.GET("/users", deps.UserHandler.GetUsers)
	public.GET("/users/:username", deps.UserHandler.FindUsername)

	protected := router.Group("/api")
	protected.Use(middlewares.JwtAuthMiddleware())

	protected.GET("/me", deps.UserHandler.CurrentUser)

	protected.POST("/servers", deps.ServerHandler.CreateServer)
	protected.GET("/servers", deps.ServerHandler.GetServers)
	protected.GET("/servers/:id", deps.ServerHandler.GetServer)
	protected.PATCH("/servers/:id", deps.ServerHandler.UpdateServer)
	protected.DELETE("/servers/:id", deps.ServerHandler.DeleteServer)

	protected.POST("/servers/:id/members", deps.ServerHandler.AddMember)
	protected.DELETE("/servers/:id/members/:userId", deps.ServerHandler.RemoveMember)
	protected.PATCH("/servers/:id/members/:userId/role", deps.ServerHandler.ChangeMemberRole)

	// protected.POST("/servers/:serverId/channels", deps.ChannelHandler.CreateChannel)
	// protected.GET("/servers/:serverId/channels", deps.ChannelHandler.GetChannels)

	// protected.POST("/channels/:channelId/messages", deps.MessageHandler.SendMessage)
	// protected.GET("/channels/:channelId/messages", deps.MessageHandler.GetMessages)

	// Start the server
	router.Run("localhost:8080")
}
