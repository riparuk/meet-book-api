package router

import (
	"github.com/gin-gonic/gin"
	"github.com/riparuk/go-gin-starter-simple/internal/database"
	"github.com/riparuk/go-gin-starter-simple/internal/handler"
	"github.com/riparuk/go-gin-starter-simple/internal/middleware"
	"github.com/riparuk/go-gin-starter-simple/internal/repository"
)

func SetupRoutes(r *gin.Engine) {

	authRepo := repository.NewUserRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)

	authHandler := handler.NewAuthHandler(authRepo)
	userHandler := handler.NewUserHandler(userRepo)

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
	}

}
