package users

import (
	"backend/src/domains/entities"
	"time"
)

type loginResponse struct {
	Token    string            `json:"token"`
	UserInfo *userInfoResponse `json:"user_info"`
}

func newLoginResponse(token string, user *entities.User) *loginResponse {
	return &loginResponse{
		Token:    token,
		UserInfo: newUserInfoResponse(user),
	}
}

type userInfoResponse struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserInfoResponse(user *entities.User) *userInfoResponse {
	return &userInfoResponse{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
