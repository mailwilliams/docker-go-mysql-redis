package controllers

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/mailwilliams/go-ambassador/src/database"
	"github.com/mailwilliams/go-ambassador/src/models"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Products(c *fiber.Ctx) error {
	var products []models.Product

	database.DB.Find(&products)

	return c.JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Create(&product)

	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}
	product.ID = uint(id)
	database.DB.Find(&product)

	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}
	product.ID = uint(id)

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Model(&product).Updates(&product)

	go database.ClearCache("products_frontend", "products_frontend")
	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}
	product.ID = uint(id)

	database.DB.Delete(&product)

	return nil
}

func ProductsFrontend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()

	result, err := database.Cache.Get(ctx, "products_frontend").Result()
	if err != nil {
		switch err {
		case redis.Nil:
			database.DB.Find(&products)

			bytes, err := json.Marshal(products)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": err.Error(),
				})
			}

			if err := database.Cache.Set(ctx, "products_frontend", bytes, 30*time.Minute).Err(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": err.Error(),
				})
			}

		default:
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": err.Error(),
			})
		}

	} else {
		if err := json.Unmarshal([]byte(result), &products); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	return c.JSON(products)
}

func ProductsBackend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()

	result, err := database.Cache.Get(ctx, "products_backend").Result()
	if err != nil {
		switch err {
		case redis.Nil:
			database.DB.Find(&products)

			bytes, err := json.Marshal(products)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": err.Error(),
				})
			}

			if err := database.Cache.Set(ctx, "products_backend", bytes, 30*time.Minute).Err(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": err.Error(),
				})
			}

		default:
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": err.Error(),
			})
		}

	} else {
		if err := json.Unmarshal([]byte(result), &products); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	var searchedProducts []models.Product

	if s := c.Query("s"); s != "" {
		lower := strings.ToLower(s)
		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), lower) || strings.Contains(strings.ToLower(product.Description), lower) {
				searchedProducts = append(searchedProducts, product)
			}
		}
	} else {
		searchedProducts = products
	}

	if sortParam := c.Query("sort"); sortParam != "" {
		sortLower := strings.ToLower(sortParam)
		if sortLower == "asc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price < searchedProducts[j].Price
			})
		} else if sortLower == "desc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price > searchedProducts[j].Price
			})
		}
	}

	var total = len(searchedProducts)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage := 9

	var data []models.Product

	if total < page*perPage && total >= (page-1)*perPage {
		data = searchedProducts[(page-1)*perPage : total]
	} else if total >= page*perPage {
		data = searchedProducts[(page-1)*perPage : page*perPage]
	} else {
		data = []models.Product{}
	}

	return c.JSON(fiber.Map{
		"data":      data,
		"total":     total,
		"page":      page,
		"last_page": total/perPage + 1,
	})
}
