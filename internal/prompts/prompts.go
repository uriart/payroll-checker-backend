package prompts

// Gemini API prompts

const GetDataFormPdfPayroll string = `
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
