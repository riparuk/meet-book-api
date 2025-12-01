package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/riparuk/go-gin-starter-simple/internal/database"
	"github.com/riparuk/go-gin-starter-simple/internal/model"
)

var allModels = []interface{}{
	&model.User{},
}

func main() {
	database.InitPostgres()

	// Enable uuid-ossp extension
	database.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	clean := false
	for _, arg := range os.Args {
		if strings.ToLower(arg) == "--clean" {
			clean = true
			break
		}
	}

	if clean {
		fmt.Println("⚠️  Dropping all tables...")
		err := database.DB.Migrator().DropTable(allModels...)
		if err != nil {
			panic("Failed to drop tables: " + err.Error())
		}

		fmt.Println("🧹 Dropping enum types...")
		// Drop types after dropping all tables
		// database.DB.Exec("DROP TYPE IF EXISTS order_status CASCADE")

		return
	}

	// Add enum types (only if not cleaning)
	// database.DB.Exec("CREATE TYPE order_status AS ENUM ('pending', 'paid', 'cancelled')")

	// Migrate schema
	err := database.DB.AutoMigrate(allModels...)
	if err != nil {
		panic("Failed to migrate: " + err.Error())
	}

	fmt.Println("✅ Database migration completed!")
}
