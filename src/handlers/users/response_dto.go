package users

import (
	"backend/src/domains/entities"
	"backend/src/handlers/common"
)

type loginResponse struct {
	Token    string                   `json:"token"`
	UserInfo *common.UserInfoResponse `json:"user_info"`
}

func newLoginResponse(token string, user *entities.User) *loginResponse {
	return &loginResponse{
		Token:    token,
		UserInfo: common.NewUserInfoResponse(user),
	}
}

type usersListResponse []*common.UserInfoResponse

func newUsersListResponse(users []*entities.User) usersListResponse {
	res := make(usersListResponse, 0, len(users))
	for _, user := range users {
		res = append(res, common.NewUserInfoResponse(user))
	}
	return res
}
