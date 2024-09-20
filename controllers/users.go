package controllers

import (
	"api/models"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Users(ctx *gin.RouterGroup) {
    ctx.GET("", GetAllUsers)
    ctx.POST("", CreateUser)
}

func GetAllUsers(ctx *gin.Context) {
    file, err := os.Open("users.json")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    defer file.Close()

    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var users []models.User

    if err = json.NewDecoder(file).Decode(&users); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, users)
}

func CreateUser(ctx *gin.Context) {
    var input models.CreateUserDto

    err := ctx.ShouldBindJSON(&input);
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := models.CreateUserDto{
        Nickname: input.Nickname,
        Email: input.Email,
        Password: input.Password,
    }
    // - Ajouter en base de donnee si les donnees sont OK (email et nickname unique)

    ctx.JSON(http.StatusOK, user)
}
