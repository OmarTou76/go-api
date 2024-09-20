package main

import (
    "github.com/gin-gonic/gin"
    "api/controllers"
)

func main() {
    router := gin.Default()
    controllers.Users(router.Group("/users"))
    router.Run(":3000")
}
