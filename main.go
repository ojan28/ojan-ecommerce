package main

import (
	"github.com/ojan28/ojan-ecommerce/controllers"
	"github.com/ojan28/ojan-ecommerce/database"
	"github.com/ojan28/ojan-ecommerce/middleware"
	"github.com/ojan28/ojan-ecommerce/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main(){
	port := os.Getenv("PORT")
	if  port == "" {
		port = "8080"
	}

	app := controllers.NewApplication(
		database.ProductData(database.Client, "Products"), 
		database.UserData(database.Client, "Users"))
	
	router := gin.New()
	router.User(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddtoCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}