package database

import (
	"log"

	"github.com/riparuk/go-gin-starter-simple/internal/model"
)

func Seed() {
	users := []model.User{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Charlie"},
	}

	for _, user := range users {
		err := DB.Create(&user).Error
		if err != nil {
			log.Printf("seeding error: %v", err)
		}
	}
}
