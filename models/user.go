package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" validate:"required" gorm:"unique;index"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
