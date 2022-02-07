package database

import (
	"github.com/mailwilliams/go-ambassador/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	DB, err = gorm.Open(mysql.Open("root:password@tcp(db:3306)/ambassador"), &gorm.Config{})
	if err != nil {
		panic("Could not connect with the database")
	}
}

func AutoMigrate() {
	err := DB.AutoMigrate(models.User{}, models.Product{}, models.Link{}, models.Order{}, models.OrderItem{})
	if err != nil {
		panic("Could not auto migrate")
	}
}
