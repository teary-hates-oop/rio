package setup

import (
	"rio/internal/db"
	"rio/internal/handlers"
	serverRepo "rio/internal/repository/server"
	userRepo "rio/internal/repository/user"
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

	userRepository := userRepo.NewDBUserRepository()
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	serverRepository := serverRepo.NewDBServerRepository()
	serverService := service.NewServerService(serverRepository, userRepository)
	serverHandler := handlers.NewServerHandler(serverService)

	// channelRepository := channel.NewDBChannelRepository()
	// channelService := service.NewChannelService(channelRepository, serverRepository)
	// channelHandler := handlers.NewChannelHandler(channelService)

	// messageRepository := message.NewDBMessageRepository()
	// messageService := service.NewMessageService(messageRepository, channelRepository, userRepository)
	// messageHandler := handlers.NewMessageHandler(messageService)

	return &Dependencies{
		UserHandler:   userHandler,
		ServerHandler: serverHandler,
		// ChannelHandler: channelHandler,
		// MessageHandler: messageHandler,
	}
}
