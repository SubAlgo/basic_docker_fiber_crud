package controllers

import (
	"mime/multipart"
)

type CreateUserForm struct {
	Username string                `form:"username" validate:"required"`
	Password string                `form:"password" validate:"required,min=4,max=15"`
	Name     string                `form:"name" validate:"required"`
	Surname  string                `form:"surname" validate:"required"`
	Email    string                `form:"email" validate:"required,email"`
	Image    *multipart.FileHeader `form:"image" validate:"required"`
}

type updateUserForm struct {
	Name    string                `form:"name"`
	Surname string                `form:"surname"`
	Image   *multipart.FileHeader `form:"image"`
}

type deleteUserForm struct {
	Disable bool `json:"disable"`
}

type userResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Image    string `json:"image"`
}

type userPaging struct {
	Items  []userResponse `json:"items"`
	Paging *pagingResult  `json:"paging"`
}
