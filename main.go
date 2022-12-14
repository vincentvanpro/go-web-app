package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	method "github.com/bu/gin-method-override"
	gin "github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/vincentvanpro/customer-web-app/internals/database"
	"github.com/vincentvanpro/customer-web-app/internals/utils"
	"github.com/vincentvanpro/customer-web-app/internals/models"
	"github.com/vincentvanpro/customer-web-app/internals/router"
)

func renderHomePage(c *gin.Context) {
	customers := []models.Customer{}
	database.DB.Find(&customers)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"customers": &customers,
		"PageTitle": "All customers",
	})
}

func MethodOverride(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Only act on POST requests.
        if r.Method == "POST" {

            // Look in the request body and headers for a spoofed method.
            // Prefer the value in the request body if they conflict.
            method := r.PostFormValue("_method")
            if method == "" {
                method = r.Header.Get("X-HTTP-Method-Override")
            }

            // Check that the spoofed method is a valid HTTP method and
            // update the request object accordingly.
            if method == "PUT" || method == "PATCH" || method == "DELETE" {
                r.Method = method
            }
        }

        // Call the next handler in the chain.
        next.ServeHTTP(w, r)
    })
}

//Custom validation (allowed from 18 till 60 years)
var ageRestrictions validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		fmt.Println(date.Year())
		fmt.Println(today.Year())
		var age = today.Year() - date.Year();
		if age <= 0 || age < 18 || age > 60 {
			return false
		}
	}
	return true
}

func main() {
	router := gin.Default()

	//Custom validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("ageRestrictions", ageRestrictions)
	}

	router.Use(method.ProcessMethodOverride(router))
	router.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("template/*.html")

	database.Connect()
	utils.GenerateSeedDataAndSaveToDB()

	//Render template paths
	router.GET("/", renderHomePage)

	//API calls
	routes.CustomerRoute(router)
	router.Run("localhost:8080")
}
