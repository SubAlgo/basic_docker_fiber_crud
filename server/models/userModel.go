package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Name     string `gorm:"not null"`
	Surname  string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Image    string `gorm:"not null"`
	RoleID   uint
	Roles    Roles `gorm:"foreignKey:RoleID"`
}

/*
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
	Image    string `json:"image"`
	//Email   string
	//Phone   string
*/
