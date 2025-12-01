package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/riparuk/go-gin-starter-simple/docs"
	"github.com/riparuk/go-gin-starter-simple/internal/database"
	"github.com/riparuk/go-gin-starter-simple/internal/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Go Gin Starter Simple
// @version 1.0.0
// @description API documentation for Go Gin Starter Simple
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func init() {
	_ = godotenv.Load(".env") // Load file .env

	docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")                // misalnya: "localhost:8080" atau "go-gin-starter-simple-api.a.run.app"
	docs.SwaggerInfo.Schemes = []string{os.Getenv("SWAGGER_SCHEME")} // atau "http" untuk lokal
}

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func main() {
	database.InitPostgres()

	r := gin.Default()

	r.Use(CORSMiddleware())
	router.SetupRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
