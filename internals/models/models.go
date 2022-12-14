package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	FirstName   string    `json:"first_name" binding:"required,max=100"`
	LastName    string    `json:"last_name" binding:"required,max=100"`
	Email       string    `json:"email" binding:"required,email"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required,ageRestrictions"`
	Gender      string    `json:"gender" binding:"required"`
	Address     string    `json:"address" binding:"max=200"`
}
