package main

import (
	"log"
	"payroll-checker-backend/api"
	"payroll-checker-backend/internal/repository"
)

func main() {

	repository.ConnectDatabase()

	router := api.SetupRouter()

	log.Printf("Servidor iniciado en el puerto 8080")
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

}
