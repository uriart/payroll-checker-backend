package api

import (
	"payroll-checker-backend/internal/handlers"
	"payroll-checker-backend/middleware"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},        // Permitir los métodos
		AllowHeaders:     []string{"Authorization", "Content-Type"}, // Permitir los encabezados
		AllowCredentials: true,                                      // Permitir las credenciales
		MaxAge:           12 * 3600,                                 // Configurar el tiempo de caché para las preflight requests
	}))

	rootGroup := router.Group("/")
	{
		rootGroup.GET("/health", handlers.HealthCheck)
	}

	defaultGroup := router.Group("/payroll")
	{
		defaultGroup.Use(middleware.EnsureValidToken()) // Grupo protegido con JWT
		defaultGroup.POST("/upload", handlers.UploadPayroll)
	}

	return router
}
