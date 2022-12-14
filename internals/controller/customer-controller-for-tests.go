package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vincentvanpro/customer-web-app/internals/database"
	"github.com/vincentvanpro/customer-web-app/internals/models"
)

func CreateCustomerTest(c *gin.Context) {
	newCustomer := models.Customer{}

	if err := c.Bind(&newCustomer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create new entity. Please check if the input is correct.",
			"debugDescr": err.Error(),
		})
		return
	}
	database.DB.Create(&newCustomer)
	c.IndentedJSON(http.StatusCreated, &newCustomer)
}

func ModifyCustomerTest(c *gin.Context) {
	id := c.Param("id")

	var customer models.Customer

	if err := database.DB.Where("id = ?", id).First(&customer).Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Object with this id not found. Cannot modify non-existent entry.",
			"debugDescr": err.Error(),
		})
		return
	}

	if err := c.Bind(&customer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update existing entity. Please check if the input is correct.",
			"debugDescr": err.Error(),
		})
		return
	}
	database.DB.Save(&customer)

	c.IndentedJSON(http.StatusOK, &customer)
}
