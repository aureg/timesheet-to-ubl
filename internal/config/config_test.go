package config

import (
	"testing"
)

func TestResolve(t *testing.T) {
	cfg := &BillingConfig{
		Default: DefaultConfig{
			HourlyRate:       100.0,
			VATPercent:       21.0,
			PaymentTermsDays: 30,
			Currency:         "EUR",
		},
		Clients: map[string]ClientConfig{
			"ClientA": {
				VATPercent: floatPtr(6.0),
			},
		},
		Projects: map[string]ProjectConfig{
			"PROJ1": {
				HourlyRate:   floatPtr(120.0),
				InvoiceLabel: "Label Proj1",
			},
		},
		ClientProjects: map[string]ProjectConfig{
			"ClientA::PROJ1": {
				HourlyRate: floatPtr(150.0),
			},
		},
	}

	// 1. Client + Project priority
	res := cfg.Resolve("ClientA", "PROJ1")
	if res.HourlyRate != 150.0 {
		t.Errorf("Expected 150.0, got %f", res.HourlyRate)
	}
	if res.VATPercent != 6.0 {
		t.Errorf("Expected 6.0, got %f", res.VATPercent)
	}
	if res.InvoiceLabel != "Label Proj1" {
		t.Errorf("Expected 'Label Proj1', got '%s'", res.InvoiceLabel)
	}

	// 2. Client priority for VAT
	res2 := cfg.Resolve("ClientA", "OTHER")
	if res2.HourlyRate != 100.0 {
		t.Errorf("Expected default 100.0, got %f", res2.HourlyRate)
	}
	if res2.VATPercent != 6.0 {
		t.Errorf("Expected client VAT 6.0, got %f", res2.VATPercent)
	}

	// 3. Project priority for rate
	res3 := cfg.Resolve("OtherClient", "PROJ1")
	if res3.HourlyRate != 120.0 {
		t.Errorf("Expected project rate 120.0, got %f", res3.HourlyRate)
	}
}

func floatPtr(f float64) *float64 { return &f }
