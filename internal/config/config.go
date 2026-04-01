package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func GetDefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}
	return filepath.Join(home, ".timesheet2ubl", "config.yaml"), nil
}

type AddressConfig struct {
	Street     string `yaml:"street"`
	City       string `yaml:"city"`
	PostalCode string `yaml:"postal_code"`
	Country    string `yaml:"country"`
}

type SupplierConfig struct {
	Name      string        `yaml:"name"`
	CompanyID string        `yaml:"company_id"`
	Email     string        `yaml:"email"`
	Phone     string        `yaml:"phone"`
	IBAN      string        `yaml:"iban"`
	BIC       string        `yaml:"bic"`
	Address   AddressConfig `yaml:"address"`
}

type CustomerConfig struct {
	Name      string        `yaml:"name"`
	CompanyID string        `yaml:"company_id"`
	Address   AddressConfig `yaml:"address"`
}

type DefaultConfig struct {
	ConsultantName   string         `yaml:"consultant_name"`
	Currency         string         `yaml:"currency"`
	VATPercent       float64        `yaml:"vat_percent"`
	HourlyRate       float64        `yaml:"hourly_rate"`
	PaymentTermsDays int            `yaml:"payment_terms_days"`
	Supplier         SupplierConfig `yaml:"supplier"`
	InvoiceTemplate  string         `yaml:"invoice_template"`
	OutputDir        string         `yaml:"output_dir"`
}

type ClientConfig struct {
	Customer         CustomerConfig `yaml:"customer"`
	VATPercent       *float64       `yaml:"vat_percent"`
	PaymentTermsDays *int           `yaml:"payment_terms_days"`
	OrderReference   string         `yaml:"order_reference"`
	InvoiceTemplate  string         `yaml:"invoice_template"`
	OutputDir        string         `yaml:"output_dir"`
}

type ProjectConfig struct {
	HourlyRate   *float64 `yaml:"hourly_rate"`
	InvoiceLabel string   `yaml:"invoice_label"`
}

type BillingConfig struct {
	Default        DefaultConfig            `yaml:"default"`
	Clients        map[string]ClientConfig  `yaml:"clients"`
	Projects       map[string]ProjectConfig `yaml:"projects"`
	ClientProjects map[string]ProjectConfig `yaml:"client_projects"`
}

func LoadConfig(path string) (*BillingConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config BillingConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

type ResolvedConfig struct {
	ConsultantName   string
	HourlyRate       float64
	VATPercent       float64
	PaymentTermsDays int
	InvoiceLabel     string
	Currency         string
	Supplier         SupplierConfig
	Customer         CustomerConfig
	OrderReference   string
	InvoiceTemplate  string
	OutputDir        string
}

func (c *BillingConfig) GetDefaultClientName() string {
	if len(c.Clients) == 1 {
		for name := range c.Clients {
			return name
		}
	}
	return ""
}

func (c *BillingConfig) Resolve(clientName, projectCode string) ResolvedConfig {
	res := ResolvedConfig{
		ConsultantName:   c.Default.ConsultantName,
		HourlyRate:       c.Default.HourlyRate,
		VATPercent:       c.Default.VATPercent,
		PaymentTermsDays: c.Default.PaymentTermsDays,
		Currency:         c.Default.Currency,
		Supplier:         c.Default.Supplier,
		InvoiceTemplate:  c.Default.InvoiceTemplate,
		OutputDir:        c.Default.OutputDir,
	}

	// 1. Client + Project (Nom Client::CODEPROJET)
	cpKey := fmt.Sprintf("%s::%s", clientName, projectCode)
	if cp, ok := c.ClientProjects[cpKey]; ok {
		if cp.HourlyRate != nil {
			res.HourlyRate = *cp.HourlyRate
		}
		if cp.InvoiceLabel != "" {
			res.InvoiceLabel = cp.InvoiceLabel
		}
	}

	// 2. Client
	if cl, ok := c.Clients[clientName]; ok {
		res.Customer = cl.Customer
		if cl.VATPercent != nil {
			res.VATPercent = *cl.VATPercent
		}
		if cl.PaymentTermsDays != nil {
			res.PaymentTermsDays = *cl.PaymentTermsDays
		}
		if cl.OrderReference != "" {
			res.OrderReference = cl.OrderReference
		}
		if cl.InvoiceTemplate != "" {
			res.InvoiceTemplate = cl.InvoiceTemplate
		}
		if cl.OutputDir != "" {
			res.OutputDir = cl.OutputDir
		}
	}

	// 3. Project
	if p, ok := c.Projects[projectCode]; ok {
		if p.HourlyRate != nil && res.HourlyRate == c.Default.HourlyRate {
			res.HourlyRate = *p.HourlyRate
		}
		if p.InvoiceLabel != "" && res.InvoiceLabel == "" {
			res.InvoiceLabel = p.InvoiceLabel
		}
	}

	if res.InvoiceLabel == "" {
		res.InvoiceLabel = fmt.Sprintf("Prestations pour le projet %s", projectCode)
	}

	return res
}

func ResolveTemplatePath(templateName string, resolvedTemplate string) string {
	// If templateName is provided via CLI, it takes priority
	path := templateName
	if path == "" {
		path = resolvedTemplate
	}

	// If still empty, search in .timesheet2ubl
	if path == "" {
		home, _ := os.UserHomeDir()
		baseDir := filepath.Join(home, ".timesheet2ubl")

		// Try default names if nothing is specified
		for _, name := range []string{"invoice.html", "invoice.docx"} {
			p := filepath.Join(baseDir, name)
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
		return ""
	}

	// If path has an extension and exists, return it
	if filepath.Ext(path) != "" {
		if _, err := os.Stat(path); err == nil {
			return path
		}
		// If it doesn't exist as is, check in .timesheet2ubl
		home, _ := os.UserHomeDir()
		p := filepath.Join(home, ".timesheet2ubl", path)
		if _, err := os.Stat(p); err == nil {
			return p
		}
		return path // Return original if not found in home
	}

	// If no extension, search for .html or .docx in .timesheet2ubl
	home, _ := os.UserHomeDir()
	baseDir := filepath.Join(home, ".timesheet2ubl")
	for _, ext := range []string{".html", ".docx"} {
		p := filepath.Join(baseDir, templateName+ext)
		if _, err := os.Stat(p); err == nil {
			return p
		}
		// Also try current directory
		if _, err := os.Stat(templateName + ext); err == nil {
			return templateName + ext
		}
	}

	return path
}
