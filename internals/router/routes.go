package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vincentvanpro/customer-web-app/internals/controller"
)

func CustomerRoute(router *gin.Engine) {
	router.GET("/users", controller.GetCustomers)
	router.GET("/users/:id", controller.GetCustomerById)
	router.POST("/users", controller.CreateCustomer)
	router.PATCH("/users/:id", controller.ModifyCustomer)
}
