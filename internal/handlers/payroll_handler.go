package handlers

import (
	"log"
	"net/http"
	"payroll-checker-backend/internal/models"
	"payroll-checker-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// UploadPayroll maneja la solicitud para subir un archivo PDF o imagen
func StructurePayroll(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Error al obtener el archivo: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Archivo no enviado o inválido"})
		return
	}
	defer file.Close()

	if header.Header.Get("Content-Type") != "application/pdf" && header.Header.Get("Content-Type") != "image/jpeg" {
		log.Printf("Archivo con MIME type inválido: %s", header.Header.Get("Content-Type"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo debe ser un PDF o una imagen jpeg"})
		return
	}

	payroll, err := services.StructurePayrollData(file, header)
	if err != nil {
		log.Printf("Error al estructurar los datos de la nómina: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Revisa los datos de la nómina"})
		return
	}

	// Responder con éxito
	log.Printf("Archivo %s procesado correctamente", header.Filename)
	c.JSON(http.StatusOK, payroll)
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func CreatePayroll(c *gin.Context) {
	var payload models.Nomina

	// Parseamos el JSON recibido
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload.EmpleadoID = c.GetString("userID")

	err := services.CreatePayroll(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nomina created successfully"})
}

func GetUserPayrolls(c *gin.Context) {
	userId := c.GetString("userID")

	nominas, err := services.GetPayrollsByUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, nominas)
}
