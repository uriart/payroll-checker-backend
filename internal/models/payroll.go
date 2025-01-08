package models

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

// Nomina representa toda la información de una nómina combinada en una sola tabla.
type Nomina struct {
	ID                             uint    `gorm:"primaryKey"`
	EmpleadoID                     string  `json:"empleado_id"`
	PuestoTrabajo                  string  `json:"puesto_trabajo"`
	GrupoProfesional               string  `json:"grupo_profesional"`
	EmpresaNombre                  string  `json:"empresa_nombre"`
	PeriodoMes                     string  `json:"periodo_mes"`
	PeriodoDiasTrabajados          int     `json:"periodo_dias_trabajados"`
	SalarioBase                    float64 `json:"salario_base"`
	Complementos                   float64 `json:"complementos"`
	HorasExtras                    float64 `json:"horas_extras"`
	Dietas                         float64 `json:"dietas"`
	PagasExtraordinarias           float64 `json:"pagas_extraordinarias"`
	TotalDevengado                 float64 `json:"total_devengado"`
	ContingenciasComunes           float64 `json:"contingencias_comunes"`
	FormacionProfesional           float64 `json:"formacion_profesional"`
	Desempleo                      float64 `json:"desempleo"`
	HorasExtrasSeguridadSocial     float64 `json:"horas_extras_seguridad_social"`
	IRPF                           float64 `json:"irpf"`
	OtrosDescuentos                float64 `json:"otros_descuentos"`
	TotalDeducido                  float64 `json:"total_deducido"`
	BaseIRPF                       float64 `json:"base_irpf"`
	BaseContingenciasComunes       float64 `json:"base_contingencias_comunes"`
	BaseContingenciasProfesionales float64 `json:"base_contingencias_profesionales"`
	LiquidoPercibido               float64 `json:"liquido_percibido"`
}

func NewPayroll(jsonPayroll *genai.GenerateContentResponse) *Nomina {
	var payroll *Nomina
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
