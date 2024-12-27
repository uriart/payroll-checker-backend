package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"payroll-checker-backend/internal/models"
	"payroll-checker-backend/internal/prompts"

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

func StructurePayrollData(file multipart.File) (*models.Payroll, error) {
	options := genai.UploadFileOptions{
		DisplayName: "payroll.pdf",
		MIMEType:    "application/pdf",
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
