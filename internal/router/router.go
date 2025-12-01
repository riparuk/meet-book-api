package router

import (
	"github.com/gin-gonic/gin"
	"github.com/riparuk/meet-book-api/internal/database"
	"github.com/riparuk/meet-book-api/internal/handler"
	"github.com/riparuk/meet-book-api/internal/middleware"
	"github.com/riparuk/meet-book-api/internal/repository"
)

func SetupRoutes(r *gin.Engine) {

	authRepo := repository.NewUserRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	roomRepo := repository.NewRoomRepository(database.DB)

	authHandler := handler.NewAuthHandler(authRepo)
	userHandler := handler.NewUserHandler(userRepo)
	roomHandler := handler.NewRoomHandler(roomRepo)

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.POST("", userHandler.CreateUser)
		}

		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}

		me := api.Group("/me")
		me.Use(middleware.JWTAuthMiddleware())
		{
			me.GET("", userHandler.Profile)
		}

		rooms := api.Group("/rooms")
		{
			rooms.GET("", roomHandler.GetRooms)
			rooms.POST("", roomHandler.CreateRoom)
			rooms.GET("/:id", roomHandler.GetRoom)
			rooms.PUT("/:id", roomHandler.UpdateRoom)
			rooms.DELETE("/:id", roomHandler.DeleteRoom)
		}
	}

}
