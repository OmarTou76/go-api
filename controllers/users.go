package controllers

import (
	"api/models"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Users(ctx *gin.RouterGroup) {
    ctx.GET("", GetAllUsers)
    ctx.POST("", CreateUser)

    ctx.GET("/:userId", GetUser)
}

func FetchUsers() ([]models.User, error) {
    file, err := os.Open("users.json")
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var users []models.User

    if err = json.NewDecoder(file).Decode(&users); err != nil {
        return nil, err
    }
    return users, nil
}

func GetAllUsers(ctx *gin.Context) {
    usersData, err := FetchUsers()
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, usersData)
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

func GetUser(ctx *gin.Context) {
    users, err := FetchUsers()
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    idStr := ctx.Params.ByName("userId")

    id, err := strconv.Atoi(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
        return
    }

    for _, user := range users {
        if user.ID == uint(id) {
            ctx.JSON(http.StatusOK, user)
            return
        }
    }
    ctx.String(http.StatusBadRequest, "User not found")
}


