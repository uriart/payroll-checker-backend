package api

import (
	"payroll-checker-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	defaultGroup := router.Group("/payroll")
	{
		defaultGroup.POST("/upload", handlers.UploadPayroll)
	}

	return router
}
