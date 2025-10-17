package databases

type createDatabaseRequestDto struct {
	Name string `json:"name" binding:"required"`
}
