package main

import (
	"api/controllers"
	"api/database"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	database.New()

	router.Use(middleware.Token())

	users := router.Group("/users")
	auth := router.Group("/auth")

	users.Use(middleware.AuthGuard())

	controllers.Auth(auth)
	controllers.Users(users)

	router.Run(":4000")
}
