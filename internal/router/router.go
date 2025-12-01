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
	bookingRepo := repository.NewBookingRepository(database.DB)

	authHandler := handler.NewAuthHandler(authRepo)
	userHandler := handler.NewUserHandler(userRepo, bookingRepo)
	roomHandler := handler.NewRoomHandler(roomRepo)
	bookingHandler := handler.NewBookingHandler(bookingRepo)

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
			me.POST("/bookings", userHandler.CreateMyBooking)
			me.GET("/bookings", userHandler.GetMyBookings)
		}

		rooms := api.Group("/rooms")
		{
			rooms.GET("", roomHandler.GetRooms)
			rooms.POST("", roomHandler.CreateRoom)
			rooms.GET("/:id", roomHandler.GetRoom)
			rooms.PUT("/:id", roomHandler.UpdateRoom)
			rooms.DELETE("/:id", roomHandler.DeleteRoom)
		}

		bookings := api.Group("/bookings")
		{
			bookings.GET("/upcoming", bookingHandler.GetUpcomingBookings)
			bookings.POST("", bookingHandler.CreateBooking)
			bookings.GET("/:id", bookingHandler.GetBooking)
			bookings.PUT("/:id", bookingHandler.UpdateBooking)
			bookings.GET("/room/:room_id", bookingHandler.GetRoomBookings)
			bookings.GET("/room/:room_id/:date", bookingHandler.GetRoomBookingsByDate)
			bookings.POST("/:id/cancel", bookingHandler.CancelBooking)
			bookings.GET("/users/:user_id", bookingHandler.GetUserBookings)
		}
	}

}
