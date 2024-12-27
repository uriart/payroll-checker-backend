package handlers

import (
	"log"
	"net/http"
	"payroll-checker-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// UploadPayroll maneja la solicitud para subir un archivo PDF
func UploadPayroll(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Error al obtener el archivo: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Archivo no enviado o inválido"})
		return
	}
	defer file.Close()

	if header.Header.Get("Content-Type") != "application/pdf" {
		log.Printf("Archivo con MIME type inválido: %s", header.Header.Get("Content-Type"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo debe ser un PDF"})
		return
	}

	payroll, err := services.StructurePayrollData(file)
	if err != nil {
		log.Printf("Error al estructurar los datos de la nómina: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Revisa los datos de la nómina"})
		return
	}

	// Responder con éxito
	log.Printf("Archivo %s procesado correctamente", header.Filename)
	c.JSON(http.StatusOK, payroll)
}
