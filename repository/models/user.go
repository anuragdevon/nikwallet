package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	EmailID  string `gorm:"unique"`
	Password string
}
