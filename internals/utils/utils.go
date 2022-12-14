package utils

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/vincentvanpro/customer-web-app/internals/database"
	"github.com/vincentvanpro/customer-web-app/internals/models"
)

func GenerateSeedDataAndSaveToDB() {
	customers := []models.Customer{}
	database.DB.Find(&customers)
	if len(customers) < 10 {
			for i := 0; i < 100; i++ {
				if rand.Intn(2) == 0 {
					database.DB.Create(&models.Customer{
						FirstName:      faker.FirstName(),
						LastName: 		faker.LastName(), 
						Email:      	strings.ToLower(faker.Email()),
						DateOfBirth: 	time.Date(2002, time.Month(rand.Intn(12)+1), 21, 1, 10, 30, 0, time.UTC),
						Gender:   		"Male",
						Address: 		faker.Timezone(),
					})
				} else {
					database.DB.Create(&models.Customer{
						FirstName:      faker.FirstName(),
						LastName: 		faker.LastName(), 
						Email:      	strings.ToLower(faker.Email()),
						DateOfBirth: 	time.Date(2002, time.Month(rand.Intn(12)+1), 21, 1, 10, 30, 0, time.UTC),
						Gender:   		"Female",
						Address: 		faker.Timezone(),
					})
				}
			}
	}

}