package invoice

import (
	"fmt"
	"math"
	"time"

	"ublcli/internal/config"
	"ublcli/internal/domain"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Calculate(entries []domain.TimesheetEntry, cfg *config.BillingConfig, invoiceNumber string) (*domain.Invoice, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("no timesheet entries provided")
	}

	// Determine client name from config (ignore excel's client name column as requested)
	clientName := cfg.GetDefaultClientName()
	if clientName == "" && len(entries) > 0 {
		// Fallback to first entry's client if config has multiple or zero clients
		clientName = entries[0].ClientName
	}

	var minDate, maxDate time.Time
	minDate = entries[0].Date
	maxDate = entries[0].Date

	var totalHours float64
	for _, e := range entries {
		if e.Date.Before(minDate) {
			minDate = e.Date
		}
		if e.Date.After(maxDate) {
			maxDate = e.Date
		}
		totalHours += e.Hours
	}

	// Period must start on the 1st of the month and end on the last day of the month
	periodStart := time.Date(minDate.Year(), minDate.Month(), 1, 0, 0, 0, 0, minDate.Location())
	periodEnd := time.Date(maxDate.Year(), maxDate.Month()+1, 0, 0, 0, 0, 0, maxDate.Location())

	invoice := &domain.Invoice{
		Number:    invoiceNumber,
		IssueDate: maxDate,
		Period: domain.Period{
			Start: periodStart,
			End:   periodEnd,
		},
		Currency: cfg.Default.Currency,
	}

	// Calculate Belgian Structured Communication (VCS)
	// Format: Year (4 digits) + Invoice Number (4 digits) + "00" + Checksum (2 digits)
	// Example: 2026 0001 00 27
	year := maxDate.Year()
	baseStr := fmt.Sprintf("%04d%s00", year, invoiceNumber)
	// Calculate checksum: base % 97. If 0, checksum is 97.
	var baseInt int64
	fmt.Sscanf(baseStr, "%d", &baseInt)
	checksum := baseInt % 97
	if checksum == 0 {
		checksum = 97
	}
	invoice.StructuredCommunication = fmt.Sprintf("+++%04d/%s/00%02d+++", year, invoiceNumber, checksum)

	// Calculate total days (8 hours per day)
	totalDays := totalHours / 8.0

	// Use default or first project for resolution to get the rate and label
	// Since there is only one line, we pick the first project code from entries or a generic one
	projectCode := entries[0].Project
	res := cfg.Resolve(clientName, projectCode)

	netAmount := s.round(totalDays * res.HourlyRate)
	taxAmount := s.round(netAmount * (res.VATPercent / 100.0))
	grossAmount := netAmount + taxAmount

	// Calculate Invoice Label (Description)
	// Pattern: "Période du 01/03/2026 au 31/03/2026 – Bon de commande n° 4110023514"
	// or "Période du 01/03/2026 au 31/03/2026"
	description := fmt.Sprintf("Période du %s au %s",
		invoice.Period.Start.Format("02/01/2006"),
		invoice.Period.End.Format("02/01/2006"),
	)
	if res.OrderReference != "" {
		description = fmt.Sprintf("%s – Bon de commande n° %s", description, res.OrderReference)
	}

	line := domain.InvoiceLine{
		ProjectCode: projectCode,
		Description: description,
		Quantity:    totalDays,
		UnitPrice:   res.HourlyRate,
		TaxPercent:  res.VATPercent,
		NetAmount:   netAmount,
		TaxAmount:   taxAmount,
		GrossAmount: grossAmount,
	}

	invoice.Lines = append(invoice.Lines, line)
	invoice.TotalHours = totalHours
	invoice.Subtotal = s.round(netAmount)
	invoice.VATAmount = s.round(taxAmount)
	invoice.TotalAmount = s.round(grossAmount)

	// Party info
	invoice.Customer = domain.Party{
		Name:      res.Customer.Name,
		CompanyID: res.Customer.CompanyID,
		Address: domain.Address{
			Street:     res.Customer.Address.Street,
			City:       res.Customer.Address.City,
			PostalCode: res.Customer.Address.PostalCode,
			Country:    res.Customer.Address.Country,
		},
	}
	invoice.Supplier = domain.Party{
		Name:      res.Supplier.Name,
		CompanyID: res.Supplier.CompanyID,
		Email:     res.Supplier.Email,
		Phone:     res.Supplier.Phone,
		IBAN:      res.Supplier.IBAN,
		BIC:       res.Supplier.BIC,
		Address: domain.Address{
			Street:     res.Supplier.Address.Street,
			City:       res.Supplier.Address.City,
			PostalCode: res.Supplier.Address.PostalCode,
			Country:    res.Supplier.Address.Country,
		},
	}
	invoice.Currency = res.Currency
	invoice.DueDate = maxDate.AddDate(0, 0, res.PaymentTermsDays)
	invoice.OrderReference = res.OrderReference
	invoice.InvoiceTemplate = res.InvoiceTemplate
	invoice.OutputDir = res.OutputDir

	return invoice, nil
}

func (s *Service) round(val float64) float64 {
	return math.Round(val*100) / 100
}
