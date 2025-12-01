package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/riparuk/go-gin-starter-simple/internal/model"
)

func Seed() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("seeding error: %v", err)
		return
	}
	users := []model.User{
		{Name: "Alice", Email: "alice@example.com", Password: string(hashedPassword)},
		{Name: "Bob", Email: "bob@example.com", Password: string(hashedPassword)},
		{Name: "Charlie", Email: "charlie@example.com", Password: string(hashedPassword)},
	}

	for _, user := range users {
		err := DB.Create(&user).Error
		if err != nil {
			log.Printf("seeding error: %v", err)
		}
	}
}
