package word

import (
	"fmt"
	"strings"
	"time"

	"ublcli/internal/domain"

	"github.com/nguyenthenguyen/docx"
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
	rDoc, err := docx.ReadDocxFile(r.templatePath)
	if err != nil {
		return fmt.Errorf("failed to read docx template: %w", err)
	}
	defer rDoc.Close()

	doc := rDoc.Editable()

	// Prepare data
	data := map[string]string{
		"InvoiceNumber":           invoice.Number,
		"InvoiceDate":             invoice.IssueDate.Format("02/01/2006"),
		"DueDate":                 invoice.DueDate.Format("02/01/2006"),
		"ClientName":              invoice.Customer.Name,
		"PeriodStart":             invoice.Period.Start.Format("02/01/2006"),
		"PeriodEnd":               invoice.Period.End.Format("02/01/2006"),
		"TotalHours":              fmt.Sprintf("%.2f", invoice.TotalHours),
		"Subtotal":                fmt.Sprintf("%.2f", invoice.Subtotal),
		"VATAmount":               fmt.Sprintf("%.2f", invoice.VATAmount),
		"TotalAmount":             fmt.Sprintf("%.2f", invoice.TotalAmount),
		"Currency":                invoice.Currency,
		"SupplierName":            invoice.Supplier.Name,
		"CustomerName":            invoice.Customer.Name,
		"OrderReference":          invoice.OrderReference,
		"StructuredCommunication": invoice.StructuredCommunication,
		"Now":                     time.Now().Format("02/01/2006 15:04"),
	}

	// Simple replacements
	for k, v := range data {
		doc.Replace(fmt.Sprintf("{{%s}}", k), v, -1)
	}

	// For lines, it's tricky in Word with a simple library.
	// Usually one would need table row duplication.
	// But as a first step, we provide a summary string for {{Lines}}
	var linesSummary strings.Builder
	for _, l := range invoice.Lines {
		linesSummary.WriteString(fmt.Sprintf("%s: %.2f h x %.2f = %.2f %s\n",
			l.Description, l.Quantity, l.UnitPrice, l.NetAmount, invoice.Currency))
	}
	doc.Replace("{{Lines}}", linesSummary.String(), -1)

	if err := doc.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("failed to write docx: %w", err)
	}

	return nil
}
