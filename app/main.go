package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"payroll-checker-backend/structs"
	"strings"

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

func structureData(path string) (string, error) {

	//imageData, err := os.ReadFile(path)
	//if err != nil {
	//	log.Fatalf("Error reading image file: %v", err)
	//}

	prompt := `
		Extract the following fields from the provided payroll:
		1. Total bruto.
		2. Total retenciones.
		3. Salario neto.
		4. Fecha de emisión.
				
		Return the JSON data only. JSON format response:
		{
		"TotalBruto": Number,
		"TotalRetenciones": Number,
		"SalarioNeto": Number,
		"FechaEmision": String
		}`

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	options := genai.UploadFileOptions{
		DisplayName: path,
		MIMEType:    "application/pdf",
	}
	fileData, err := client.UploadFile(ctx, "", file, &options)
	if err != nil {
		log.Fatalf("Error uploading file: %v", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.FileData{URI: fileData.URI, MIMEType: fileData.MIMEType}, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	var payroll structs.Payroll
	part := resp.Candidates[0].Content.Parts[0]
	if text, ok := part.(genai.Text); ok {
		var parsedText, _ = parseAIResponse(string(text))
		err := json.Unmarshal(parsedText, &payroll)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON: %v", err)
		}
	}

	log.Printf("Nómina estructurada y formateada %v", payroll)

	return "", nil

}

func parseAIResponse(aiResponse string) ([]byte, error) {
	// Detecta el bloque de código Markdown
	if strings.HasPrefix(aiResponse, "```json") {
		// Elimina las líneas de inicio y fin del bloque de código
		start := strings.Index(aiResponse, "```json") + len("```json\n")
		end := strings.LastIndex(aiResponse, "```")
		if start < end {
			aiResponse = aiResponse[start:end]
		}
	}

	// Limpia espacios innecesarios
	aiResponse = strings.TrimSpace(aiResponse)

	// Retorna la respuesta limpia como bytes
	return []byte(aiResponse), nil
}

func main() {
	filePath := "payroll_sample.pdf"
	_, err := structureData(filePath)
	if err != nil {
		log.Fatal(err)
	}

}
