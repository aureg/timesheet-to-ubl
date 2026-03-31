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

	// 4. Template resolution
	cfg.Default.InvoiceTemplate = "default-tpl"
	cfg.Clients["ClientA"] = ClientConfig{
		InvoiceTemplate: "client-tpl",
	}

	res4 := cfg.Resolve("ClientA", "PROJ1")
	if res4.InvoiceTemplate != "client-tpl" {
		t.Errorf("Expected 'client-tpl', got '%s'", res4.InvoiceTemplate)
	}

	res5 := cfg.Resolve("OtherClient", "PROJ1")
	if res5.InvoiceTemplate != "default-tpl" {
		t.Errorf("Expected 'default-tpl', got '%s'", res5.InvoiceTemplate)
	}

	// 5. OutputDir resolution
	cfg.Default.OutputDir = "./global-out"
	cfg.Clients["ClientA"] = ClientConfig{
		OutputDir: "./client-out",
	}

	res6 := cfg.Resolve("ClientA", "PROJ1")
	if res6.OutputDir != "./client-out" {
		t.Errorf("Expected './client-out', got '%s'", res6.OutputDir)
	}

	res7 := cfg.Resolve("OtherClient", "PROJ1")
	if res7.OutputDir != "./global-out" {
		t.Errorf("Expected './global-out', got '%s'", res7.OutputDir)
	}
}

func TestResolveTemplatePath(t *testing.T) {
	// Difficult to test exactly because it depends on os.UserHomeDir and existence of files
	// but we can check some logic

	// Ext present and exists locally (current dir)
	path := ResolveTemplatePath("invoice.html", "")
	// If invoice.html exists in current dir, it should return it as is or relative
	if path == "" {
		t.Errorf("Expected some path for invoice.html")
	}

	// No extension -> should try .html or .docx
	path2 := ResolveTemplatePath("invoice", "")
	if path2 == "" {
		t.Errorf("Expected path for 'invoice'")
	}
}

func floatPtr(f float64) *float64 { return &f }
