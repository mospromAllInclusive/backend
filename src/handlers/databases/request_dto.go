package databases

type createDatabaseRequestDto struct {
	Name string `json:"name" binding:"required"`
}

type setRoleRequestDto struct {
	UserID int64  `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=reader writer admin"`
}

type deleteUserRequestDto struct {
	UserID int64 `json:"user_id" binding:"required"`
}
