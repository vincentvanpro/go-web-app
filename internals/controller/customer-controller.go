package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vincentvanpro/customer-web-app/internals/database"
	"github.com/vincentvanpro/customer-web-app/internals/models"
)


func GetCustomers(c *gin.Context) {
	customers := []models.Customer{}
	database.DB.Find(&customers)
	c.IndentedJSON(http.StatusOK, &customers)
}

func GetCustomerById(c *gin.Context) {
	id := c.Param("id")

	var customer models.Customer
	database.DB.Where("id = ?", id).First(&customer)

	if strconv.FormatUint(uint64(customer.ID), 10) != id {
		c.JSON(http.StatusNotFound, "No customer with such ID in database")
		return
	}

	c.JSON(http.StatusOK, customer)
}

func CreateCustomer(c *gin.Context) {
	newCustomer := models.Customer{}

	if err := c.Bind(&newCustomer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create new entity. Please check if the input is correct.",
			"debugDescr": err.Error(),
		})
		return
	}
	database.DB.Create(&newCustomer)
	//c.IndentedJSON(http.StatusCreated, &newCustomer) //UNCOMMENT FOR TESTING
	c.Redirect(http.StatusFound, "/")
}

func ModifyCustomer(c *gin.Context) {
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

	// c.IndentedJSON(http.StatusOK, &customer) //UNCOMMENT FOR TESTING
	c.Redirect(http.StatusFound, "/")
}
