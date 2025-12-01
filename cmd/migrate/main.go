package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/riparuk/meet-book-api/internal/database"
	"github.com/riparuk/meet-book-api/internal/model"
)

var allModels = []interface{}{
	&model.User{},
	&model.Room{},
	&model.Booking{},
}

func main() {
	database.InitPostgres()

	// Test database connection
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	// Check if database is accessible
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}

	// Get database name for logging
	var dbName string
	err = database.DB.Raw("SELECT current_database()").Scan(&dbName).Error
	if err != nil {
		log.Fatalf("❌ Failed to get database name: %v", err)
	}
	fmt.Printf("🔍 Connected to database: %s\n", dbName)

	// Check if extension exists
	var extExists bool
	err = database.DB.Raw("SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp'").Scan(&extExists).Error
	if err != nil {
		log.Fatalf("❌ Failed to check for uuid-ossp extension: %v", err)
	}

	if !extExists {
		fmt.Println("🔧 Creating uuid-ossp extension...")
		err = database.DB.Exec("CREATE EXTENSION \"uuid-ossp\"").Error
		if err != nil {
			log.Fatalf("❌ Failed to create uuid-ossp extension: %v", err)
		}
	}

	clean := false
	for _, arg := range os.Args {
		if strings.ToLower(arg) == "--clean" {
			clean = true
			break
		}
	}

	if clean {
		fmt.Println("⚠️  Dropping all tables...")
		err = database.DB.Migrator().DropTable(allModels...)
		if err != nil {
			log.Fatalf("❌ Failed to drop tables: %v", err)
		}
		fmt.Println("✅ Successfully dropped all tables")
	}

	// Migrate schema
	fmt.Println("🔄 Running database migrations...")
	err = database.DB.AutoMigrate(allModels...)
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	// Verify tables were created
	var tables []string
	err = database.DB.Raw(
		`SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'`).
		Pluck("table_name", &tables).Error

	if err != nil {
		log.Fatalf("❌ Failed to list tables: %v", err)
	}

	fmt.Println("\n📊 Database tables:")
	if len(tables) == 0 {
		fmt.Println("No tables found in the database!")
	} else {
		for _, table := range tables {
			fmt.Printf("- %s\n", table)
		}
	}

	fmt.Println("\n✅ Database migration completed!")
}
