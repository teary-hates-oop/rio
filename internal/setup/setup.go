package setup

import (
	"rio/internal/db"
	"rio/internal/handlers"
	server "rio/internal/repository/server" // alias: server
	user "rio/internal/repository/user"     // alias: user
	"rio/internal/service"
)

type Dependencies struct {
	UserHandler   *handlers.UserHandler
	ServerHandler *handlers.ServerHandler
	// ChannelHandler *handlers.ChannelHandler
	// MessageHandler *handlers.MessageHandler
}

func Setup() *Dependencies {
	db.ConnectDataBase()

	// User layer
	userRepo := user.NewDBUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Server layer
	serverRepo := server.NewDBServerRepository()
	serverService := service.NewServerService(serverRepo, userRepo) // pass userRepo for membership checks
	serverHandler := handlers.NewServerHandler(serverService)

	// Channel layer (uncomment when ready)
	// channelRepo := channel.NewDBChannelRepository()
	// channelService := service.NewChannelService(channelRepo, serverRepo)
	// channelHandler := handlers.NewChannelHandler(channelService)

	// Message layer (uncomment when ready)
	// messageRepo := message.NewDBMessageRepository()
	// messageService := service.NewMessageService(messageRepo, channelRepo, userRepo)
	// messageHandler := handlers.NewMessageHandler(messageService)

	return &Dependencies{
		UserHandler:   userHandler,
		ServerHandler: serverHandler,
		// ChannelHandler: channelHandler,
		// MessageHandler: messageHandler,
	}
}
