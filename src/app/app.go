package app

import (
	"backend/src/app/repositories"
	"backend/src/app/resources"
	"backend/src/app/services"
	"log"

	"github.com/gin-gonic/gin"
)

type App struct {
	Engine       *gin.Engine
	Resources    *resources.Resources
	Repositories *repositories.Repositories
	Services     *services.Services
}

func (a *App) Run() {
	log.Fatal(a.Engine.Run(":8080"))
}

func (a *App) Init() {
	a.initResources()
	a.initRepositories()
	a.initServices()
	a.initWebServer()
}

func (a *App) initResources() {
	a.Resources = resources.NewResources()
}

func (a *App) initRepositories() {
	a.Repositories = repositories.NewRepositories(a.Resources)
}

func (a *App) initServices() {
	a.Services = services.NewServices(a.Repositories, a.Resources)
}
