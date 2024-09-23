package controllers

import (
	"api/models"
	"api/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Users(ctx *gin.RouterGroup) {
    ctx.GET("", GetAllUsers)
    ctx.POST("", CreateUser)

    ctx.GET("/:userId", GetUserById)
    ctx.DELETE("/:userId", DeleteUser)
}

func GetAllUsers(ctx *gin.Context) {
    var users []models.User
    if err := database.DB.Find(&users).Error; err != nil {
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

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    newUser := models.User{
        Nickname: input.Nickname,
        Email: input.Email,
        Password: string(hashedPassword),
    }

    if err := database.DB.Create(&newUser).Error; err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.Status(http.StatusCreated)
}

func GetUserById(ctx *gin.Context) {
    id, err := GetUserIdToUINT(ctx.Params)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var target models.User
    if err := database.DB.First(&target, id).Error; err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, target)
}


func DeleteUser(ctx *gin.Context) {
    id, err := GetUserIdToUINT(ctx.Params)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    result := database.DB.Delete(&models.User{}, id)

    if result.Error != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    if result.RowsAffected == 0 {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    ctx.Status(http.StatusNoContent)
}

func GetUserIdToUINT(params gin.Params) (int, error) {
    idStr := params.ByName("userId")

    id, err := strconv.Atoi(idStr)
    if err != nil {
        return 0, err
    }
    return id, nil
}

