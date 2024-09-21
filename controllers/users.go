package controllers

import (
	"api/models"
	"encoding/json"
	"net/http"
	"os"
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

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user := models.CreateUserDto{
        Nickname: input.Nickname,
        Email: input.Email,
        Password: string(hashedPassword),
    }
    users, err := FetchUsers()
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var idx uint = 1
    if len(users) > 0 {
        idx = users[len(users)-1].ID + 1
    }

    newUser := models.User{
        ID: idx,
        Nickname: user.Nickname,
        Email: user.Email,
        Password: user.Password,
    }

    users = append(users, newUser)
    if err := UpdateFile("users.json", users); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, newUser)
}

func GetUserById(ctx *gin.Context) {
    id, err := GetUserId(ctx.Params)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    users, err := FetchUsers()
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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


func DeleteUser(ctx *gin.Context) {
    id, err := GetUserId(ctx.Params)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    users, err := FetchUsers()
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    var target *models.User
    var indexToRemove int
    for i, user := range users {
        if user.ID == uint(id) {
            target = &user
            indexToRemove = i
            break
        }
    }

    if (target == nil) {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    users = append(users[:indexToRemove], users[indexToRemove + 1:]...)
    if err := UpdateFile("users.json", users); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    ctx.String(http.StatusOK, "User deleted")
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

func UpdateFile(name string, data interface{}) error {
    fileUpdated, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer fileUpdated.Close()
    bytes, err := json.Marshal(data)
    if err != nil {
        return err
    }
    _, err = fileUpdated.Write(bytes)
    if err != nil {
        return err
    }
    return nil
}

func GetUserId(params gin.Params) (int, error) {
    idStr := params.ByName("userId")

    id, err := strconv.Atoi(idStr)
    if err != nil {
        return 0, err
    }
    return id, nil
}

