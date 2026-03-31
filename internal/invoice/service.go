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

func (s *Service) Calculate(entries []domain.TimesheetEntry, cfg *config.BillingConfig) (*domain.Invoice, error) {
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

	invoice := &domain.Invoice{
		Number:    fmt.Sprintf("INV-%s", time.Now().Format("20060102-150405")), // Default number
		IssueDate: maxDate,
		Period: domain.Period{
			Start: minDate,
			End:   maxDate,
		},
		Currency: cfg.Default.Currency,
	}

	// Calculate total days (8 hours per day)
	totalDays := totalHours / 8.0

	// Use default or first project for resolution to get the rate and label
	// Since there is only one line, we pick the first project code from entries or a generic one
	projectCode := entries[0].Project
	res := cfg.Resolve(clientName, projectCode)

	netAmount := s.round(totalDays * res.HourlyRate)
	taxAmount := s.round(netAmount * (res.VATPercent / 100.0))
	grossAmount := netAmount + taxAmount

	line := domain.InvoiceLine{
		ProjectCode: projectCode,
		Description: res.InvoiceLabel,
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

	return invoice, nil
}

func (s *Service) round(val float64) float64 {
	return math.Round(val*100) / 100
}
