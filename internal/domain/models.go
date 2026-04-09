package domain

import (
	"time"
)

type Address struct {
	Street     string
	City       string
	PostalCode string
	Country    string
}

type Party struct {
	Name      string
	CompanyID string
	Email     string
	Phone     string
	IBAN      string
	BIC       string
	Address   Address
}

type Money struct {
	Amount   float64
	Currency string
}

type TimesheetEntry struct {
	Date        time.Time
	Project     string
	ProjectDesc string
	Task        string
	TaskDesc    string
	ClientName  string
	Hours       float64
	Ticket      string
	Comment     string
	LineNumber  int
}

type InvoiceLine struct {
	Description string
	Quantity    float64
	UnitPrice   float64
	TaxPercent  float64
	NetAmount   float64
	TaxAmount   float64
	GrossAmount float64
}

type Period struct {
	Start time.Time
	End   time.Time
}

type Attachment struct {
	Filename string
	MimeType string
	Content  []byte
}

type Invoice struct {
	Number                  string
	IssueDate               time.Time
	DueDate                 time.Time
	Period                  Period
	Supplier                Party
	Customer                Party
	Lines                   []InvoiceLine
	TotalHours              float64
	Subtotal                float64
	VATAmount               float64
	TotalAmount             float64
	Currency                string
	Attachments             []Attachment
	OrderReference          string
	StructuredCommunication string
	InvoiceTemplate         string
	ExcelTemplate           string
	OutputDir               string
	ConsultantName          string
	ManagerName             string
}
