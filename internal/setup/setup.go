package setup

import (
	"rio/internal/db"
	"rio/internal/handlers"
	"rio/internal/repository"
	"rio/internal/service"
)

type Dependencies struct {
	UserHandler *handlers.UserHandler
}

func Setup() *Dependencies {
	db.ConnectDataBase()

	userRepo := repository.NewDBUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	return &Dependencies{
		UserHandler: userHandler,
	}
}
