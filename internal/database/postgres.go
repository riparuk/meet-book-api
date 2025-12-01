package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: No .env file found, using environment variables.")
	}

	// Prioritize direct connection if available
	dsn := os.Getenv("DATABASE_DIRECT_URL")
	if dsn != "" {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt: false,
		})
		if err != nil {
			log.Fatalf("failed to connect to Direct Postgres: %v", err)
		}
		fmt.Println("🟢 Connected to Direct Postgres")
		return
	}

	// Remote Postgres
	dsn = os.Getenv("DATABASE_URL")
	if dsn != "" {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt: false,
		})
		if err != nil {
			log.Fatalf("failed to connect to Remote Postgres: %v", err)
		}
		fmt.Println("🟢 Connected to Remote Postgres")
		return
	}

	fmt.Println("🟢 Connected to Postgres")
}
