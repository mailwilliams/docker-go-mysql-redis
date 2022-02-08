package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/mailwilliams/go-ambassador/src/database"
	"github.com/mailwilliams/go-ambassador/src/models"
)

func main() {
	database.Connect()
	database.SetUpRedis()

	ctx := context.Background()

	var users []models.User

	database.DB.Find(&users, models.User{IsAmbassador: true})

	for _, user := range users {
		ambassador := models.Ambassador(user)
		ambassador.CalculateRevenue(database.DB)

		database.Cache.ZAdd(ctx, "rankings", &redis.Z{Score: *ambassador.Revenue, Member: user.Name()})
	}
}
