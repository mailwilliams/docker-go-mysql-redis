package main

import (
	"github.com/bxcodec/faker/v3"
	"github.com/mailwilliams/go-ambassador/src/database"
	"github.com/mailwilliams/go-ambassador/src/models"
)

func main() {
	database.Connect()
	for i := 0; i < 30; i++ {
		ambassador := models.User{
			FirstName:    faker.FirstName(),
			LastName:     faker.LastName(),
			Email:        faker.Email(),
			IsAmbassador: true,
		}

		ambassador.Password = []byte("1234")

		database.DB.Create(&ambassador)
	}
}
