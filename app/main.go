package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
	"google.golang.org/api/option"
)

var pool pdfium.Pool
var instance pdfium.Pdfium
var client *genai.Client
var ctx context.Context

func init() {

	var err error
	// Init the PDFium library and return the instance to open documents.
	pool, err = webassembly.Init(webassembly.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
	})
	if err != nil {
		log.Fatal(err)
	}

	instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}

	ctx = context.Background()
	client, err = genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
}

func textractPDF(filePath string) (string, error) {

	// Load the PDF file into a byte array.
	pdfBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Open the PDF using PDFium (and claim a worker)
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return "", err
	}

	// Always close the document, this will release its resources.
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	pages, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		return "", err
	}

	var allPagesText string
	for i := 0; i < pages.PageCount; i++ {
		pageText, _ := instance.GetPageText(&requests.GetPageText{
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    i,
				},
			},
		})
		allPagesText += pageText.Text
	}

	log.Printf("PDF con %d páginas", pages.PageCount)

	// TODO Extract text from image with tesseract-ocr/tesseract
	if false {
		_, err = instance.RenderToFile(&requests.RenderToFile{
			OutputTarget:   "file",
			TargetFilePath: "/renders/output_image.jpg",
			OutputFormat:   "jpg",
			RenderPageInDPI: &requests.RenderPageInDPI{
				DPI: 300,
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: doc.Document,
						Index:    0,
					},
				},
			},
		})
		if err != nil {
			log.Printf("Error RenderToFile %v", err)
		}
	}

	instance.Close()

	return allPagesText, nil

}

func structureData(text string) (string, error) {

	prompt := fmt.Sprintf(`
		Extract the following fields from the provided payroll:
		1. Total gross.
		2. Total deductions.
		3. Net salary.
		4. Date of issue.
		
		Payroll text:
		%s
		
		JSON format response:
		{
		“total_gross": "",
		“total_deductions": "",
		“net_salary": "",
		“issue_date": ””
		}`, text)

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			fmt.Println(txt)
		}
	}

	//log.Printf("Respuesta del modelo Gemini: %v", resp.Candidates[0].Content)

	return prompt, nil

}

func main() {
	filePath := "payroll_sample.pdf"
	text, err := textractPDF(filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Texto extraido del pdf: %v", text)

	structureData(text)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	//log.Printf("Datos de la nómina estructurados: %v", structuredData)

}
