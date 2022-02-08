package controllers

import (
	"github.com/bxcodec/faker/v3"
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

type CreateLinkRequest struct {
	Products []int
}

func CreateLink(c *fiber.Ctx) error {
	var request CreateLinkRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	id, _ := middlewares.GetUserID(c)

	link := models.Link{
		UserID: id,
		Code:   faker.Username(),
	}

	for _, productID := range request.Products {
		product := models.Product{}
		product.ID = uint(productID)
		link.Products = append(link.Products, product)
	}

	database.DB.Create(&link)

	return c.JSON(link)
}

func Stats(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserID(c)

	var links []models.Link

	database.DB.Find(&links, models.Link{UserID: id})

	var result []interface{}

	var orders []models.Order

	for _, link := range links {
		database.DB.Preload("OrderItems").Find(&orders, &models.Order{Code: link.Code, Complete: true})

		revenue := 0.0
		for _, order := range orders {
			revenue += order.GetTotal()
		}

		result = append(result, fiber.Map{"code": link.Code, "count": len(orders), "revenue": revenue})
	}

	return c.JSON(result)
}
