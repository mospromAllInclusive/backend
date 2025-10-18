package app

import (
	"backend/src/handlers"
	"backend/src/handlers/changelog"
	"backend/src/handlers/databases"
	"backend/src/handlers/tables"
	"backend/src/handlers/users"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (a *App) initWebServer() {
	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}

	r.Use(cors.New(corsConfig))

	//r.MaxMultipartMemory = MaxUploadSize

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	authMiddleware := a.Services.AuthService.JWTAuthMiddleware()
	for _, handler := range a.initHandlers() {
		h := make([]gin.HandlerFunc, 0)
		if handler.AuthRequired() {
			h = append(h, authMiddleware)
		}
		h = append(h, handler.Handle)

		r.Handle(handler.Method(), handler.Path(), h...)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	a.Engine = r
}

func (a *App) initHandlers() []handlers.IHandler {
	res := make([]handlers.IHandler, 0)

	res = append(res, tables.NewHandlers(a.Services.TablesService, a.Services.DatabasesService, a.Services.FileService)...)
	res = append(res, databases.NewHandlers(a.Services.TablesService, a.Services.DatabasesService, a.Services.UsersService)...)
	res = append(res, users.NewHandlers(a.Services.AuthService, a.Services.UsersService)...)
	res = append(res, changelog.NewHandlers(a.Services.ChangelogService, a.Services.TablesService, a.Services.DatabasesService)...)

	return res
}
