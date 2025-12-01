package main

import (
	"fmt"

	"github.com/riparuk/meet-book-api/internal/database"
)

func main() {
	database.InitPostgres()
	database.Seed()

	fmt.Println("Seeding completed")
}
