package htmltmpl

import (
	"fmt"
	"html/template"
	"os"

	"ublcli/internal/domain"
)

type Renderer struct {
	templatePath string
}

func NewRenderer(templatePath string) *Renderer {
	return &Renderer{
		templatePath: templatePath,
	}
}

func (r *Renderer) Render(invoice *domain.Invoice, outputPath string) error {
	tmpl, err := template.ParseFiles(r.templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	data := map[string]interface{}{
		"InvoiceNumber": invoice.Number,
		"InvoiceDate":   invoice.IssueDate.Format("02/01/2006"),
		"DueDate":       invoice.DueDate.Format("02/01/2006"),
		"ClientName":    invoice.Customer.Name,
		"PeriodStart":   invoice.Period.Start.Format("02/01/2006"),
		"PeriodEnd":     invoice.Period.End.Format("02/01/2006"),
		"Lines":         invoice.Lines,
		"TotalHours":    fmt.Sprintf("%.2f", invoice.TotalHours),
		"Subtotal":      fmt.Sprintf("%.2f", invoice.Subtotal),
		"VATAmount":     fmt.Sprintf("%.2f", invoice.VATAmount),
		"TotalAmount":   fmt.Sprintf("%.2f", invoice.TotalAmount),
		"Currency":      invoice.Currency,
		"Supplier":      invoice.Supplier,
		"Customer":      invoice.Customer,
	}

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
