package models

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

type Payroll struct {
	TotalBruto       float32
	TotalRetenciones float32
	SalarioNeto      float32
	FechaEmision     string
}

func NewPayroll(jsonPayroll *genai.GenerateContentResponse) *Payroll {
	var payroll *Payroll
	part := jsonPayroll.Candidates[0].Content.Parts[0]
	if text, ok := part.(genai.Text); ok {
		var parsedText, _ = parseAIResponse(string(text))
		err := json.Unmarshal(parsedText, &payroll)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON: %v", err)
		}
	}
	return payroll
}

func parseAIResponse(aiResponse string) ([]byte, error) {
	if strings.HasPrefix(aiResponse, "```json") {
		start := strings.Index(aiResponse, "```json") + len("```json\n")
		end := strings.LastIndex(aiResponse, "```")
		if start < end {
			aiResponse = aiResponse[start:end]
		}
	}
	aiResponse = strings.TrimSpace(aiResponse)
	return []byte(aiResponse), nil
}
