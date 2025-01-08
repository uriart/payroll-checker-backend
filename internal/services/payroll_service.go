package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"payroll-checker-backend/internal/models"
	"payroll-checker-backend/internal/prompts"
	"payroll-checker-backend/internal/repository"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var client *genai.Client
var ctx context.Context

func init() {
	var err error
	ctx = context.Background()
	client, err = genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
}

func StructurePayrollData(file multipart.File, header *multipart.FileHeader) (*models.Nomina, error) {
	options := genai.UploadFileOptions{
		DisplayName: header.Filename,
		MIMEType:    header.Header.Get("Content-Type"),
	}
	fileData, err := client.UploadFile(ctx, "", file, &options)
	if err != nil {
		return nil, fmt.Errorf("error uploading file: %v", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.FileData{URI: fileData.URI, MIMEType: fileData.MIMEType}, genai.Text(prompts.GetDataFormPdfPayroll))
	if err != nil {
		return nil, err
	}

	return models.NewPayroll(resp), nil
}

func CreatePayroll(payload models.Nomina) error {
	db := repository.DB

	if err := db.Create(&payload).Error; err != nil {
		return fmt.Errorf("falied to create NominaCompleta")
	}

	return nil
}

func GetPayrollsByUser(userId string) ([]models.Nomina, error) {
	var nominas []models.Nomina
	db := repository.DB

	result := db.Where("empleado_id = ?", userId).Find(&nominas)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting payrolls: %v", result.Error)
	}

	return nominas, nil
}
