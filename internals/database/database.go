package database

import (
	"fmt"

	"github.com/vincentvanpro/customer-web-app/internals/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	fmt.Println("Connected")
	db, err := gorm.Open(postgres.Open("postgres://pg:pass@localhost:5432/customers"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Customer{})
	DB=db
}

func GetDB() *gorm.DB {
	return DB
}

func ClearTable() {
	DB.Exec("DELETE FROM customers")
	DB.Exec("ALTER SEQUENCE customers_id_seq RESTART WITH 1")
}
