package main

import (
	"fmt"

	"github.com/riparuk/go-gin-starter-simple/internal/database"
)

func main() {
	database.InitPostgres()
	database.Seed()

	fmt.Println("Seeding completed")
}
