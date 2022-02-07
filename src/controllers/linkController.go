package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mailwilliams/go-ambassador/src/database"
	"github.com/mailwilliams/go-ambassador/src/middlewares"
	"github.com/mailwilliams/go-ambassador/src/models"
)

func Link(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserID(c)

	var links []models.Link
	database.DB.Find("user_id = ?", id).Find(&links)

	for i, link := range links {
		var orders []models.Order

		database.DB.Where("code = ? and complete = true", link.Code).Find(&orders)

		links[i].Orders = orders
	}
	return c.JSON(links)
}
