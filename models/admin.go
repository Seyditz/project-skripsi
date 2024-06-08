package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Id       int    `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminListResponse struct {
	Id    int    `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
