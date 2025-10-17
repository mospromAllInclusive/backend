package users

import "backend/src/domains/entities"

type loginRequestDto struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerRequestDto struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r registerRequestDto) ToUser() *entities.User {
	return &entities.User{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}
