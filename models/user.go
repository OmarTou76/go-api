package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    ID uint `json:"id" gorm:"primary_key"` 
    Nickname string `json:"nickname" gorm:"unique;not null" binding:"required,min=3"`
    Email string `json:"email" gorm:"unique;not null" binding:"required,email"`
    Password string `json:"password" gorm:"not null" binding:"required,min=6"`
    // A completer au fur et a mesure
}

type CreateUserDto struct {
    Nickname string `json:"nickname" binding:"required,min=3"`
    Email string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
