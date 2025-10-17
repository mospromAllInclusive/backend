package handlers

import "github.com/gin-gonic/gin"

type IHandler interface {
	Handle(c *gin.Context)
	Path() string
	Method() string
	AuthRequired() bool
}
