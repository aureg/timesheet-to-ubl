package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

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
	Currency         string         `yaml:"currency"`
	VATPercent       float64        `yaml:"vat_percent"`
	HourlyRate       float64        `yaml:"hourly_rate"`
	PaymentTermsDays int            `yaml:"payment_terms_days"`
	Supplier         SupplierConfig `yaml:"supplier"`
}

type ClientConfig struct {
	Customer         CustomerConfig `yaml:"customer"`
	VATPercent       *float64       `yaml:"vat_percent"`
	PaymentTermsDays *int           `yaml:"payment_terms_days"`
	OrderReference   string         `yaml:"order_reference"`
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
	HourlyRate       float64
	VATPercent       float64
	PaymentTermsDays int
	InvoiceLabel     string
	Currency         string
	Supplier         SupplierConfig
	Customer         CustomerConfig
	OrderReference   string
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
		HourlyRate:       c.Default.HourlyRate,
		VATPercent:       c.Default.VATPercent,
		PaymentTermsDays: c.Default.PaymentTermsDays,
		Currency:         c.Default.Currency,
		Supplier:         c.Default.Supplier,
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
