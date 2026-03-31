package invoice

import (
	"testing"
	"time"

	"ublcli/internal/config"
	"ublcli/internal/domain"
)

func TestCalculate(t *testing.T) {
	cfg := &config.BillingConfig{
		Default: config.DefaultConfig{
			Currency:         "EUR",
			VATPercent:       21.0,
			HourlyRate:       100.0,
			PaymentTermsDays: 30,
		},
		Clients: map[string]config.ClientConfig{
			"Client A": {
				Customer:       config.CustomerConfig{Name: "Client A Legal"},
				OrderReference: "PO-CLIENT-A",
			},
		},
	}

	entries := []domain.TimesheetEntry{
		{
			Date:       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			ClientName: "Client A",
			Project:    "P1",
			Hours:      5.0,
		},
		{
			Date:       time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			ClientName: "Client A",
			Project:    "P1",
			Hours:      3.0,
		},
		{
			Date:       time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			ClientName: "Client A",
			Project:    "P2",
			Hours:      2.0,
		},
	}

	svc := NewService()
	inv, err := svc.Calculate(entries, cfg)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if inv.TotalHours != 10.0 {
		t.Errorf("Expected 10 hours, got %f", inv.TotalHours)
	}

	if len(inv.Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(inv.Lines))
	}

	expectedDays := 10.0 / 8.0
	if inv.Lines[0].Quantity != expectedDays {
		t.Errorf("Expected %f days, got %f", expectedDays, inv.Lines[0].Quantity)
	}

	expectedSubtotal := expectedDays * 100.0
	if inv.Subtotal != expectedSubtotal {
		t.Errorf("Expected %f subtotal, got %f", expectedSubtotal, inv.Subtotal)
	}

	expectedVAT := expectedSubtotal * 0.21
	if inv.VATAmount != expectedVAT {
		t.Errorf("Expected %f VAT, got %f", expectedVAT, inv.VATAmount)
	}

	expectedTotal := expectedSubtotal + expectedVAT
	if inv.TotalAmount != expectedTotal {
		t.Errorf("Expected %f total, got %f", expectedTotal, inv.TotalAmount)
	}

	expectedStart := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if !inv.Period.Start.Equal(expectedStart) {
		t.Errorf("Expected start %v, got %v", expectedStart, inv.Period.Start)
	}

	expectedEnd := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)
	if !inv.Period.End.Equal(expectedEnd) {
		t.Errorf("Expected end %v, got %v", expectedEnd, inv.Period.End)
	}

	expectedDescription := "Période du 01/01/2023 au 31/01/2023 – Bon de commande n° PO-CLIENT-A"
	if inv.Lines[0].Description != expectedDescription {
		t.Errorf("Expected description '%s', got '%s'", expectedDescription, inv.Lines[0].Description)
	}

	if inv.OrderReference != "PO-CLIENT-A" {
		t.Errorf("Expected OrderReference 'PO-CLIENT-A', got '%s'", inv.OrderReference)
	}
}

func TestCalculate_MultipleClientsNoError(t *testing.T) {
	cfg := &config.BillingConfig{
		Clients: map[string]config.ClientConfig{
			"Client A": {
				Customer: config.CustomerConfig{Name: "Client A Legal"},
			},
		},
	}
	entries := []domain.TimesheetEntry{
		{ClientName: "Client A", Date: time.Now(), Project: "P1", Hours: 1, LineNumber: 1},
		{ClientName: "Client B", Date: time.Now(), Project: "P1", Hours: 1, LineNumber: 2},
	}

	svc := NewService()
	inv, err := svc.Calculate(entries, cfg)
	if err != nil {
		t.Fatalf("Expected no error for multiple clients from excel, got %v", err)
	}

	if inv.Customer.Name != "Client A Legal" {
		t.Errorf("Expected Customer Name 'Client A Legal', got '%s'", inv.Customer.Name)
	}
}
