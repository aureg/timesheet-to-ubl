package ubl

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"ublcli/internal/domain"
)

type Invoice struct {
	XMLName                      xml.Name                      `xml:"Invoice"`
	Xmlns                        string                        `xml:"xmlns,attr"`
	XmlnsCac                     string                        `xml:"xmlns:cac,attr"`
	XmlnsCbc                     string                        `xml:"xmlns:cbc,attr"`
	CustomizationID              string                        `xml:"cbc:CustomizationID"`
	ProfileID                    string                        `xml:"cbc:ProfileID"`
	ID                           string                        `xml:"cbc:ID"`
	IssueDate                    string                        `xml:"cbc:IssueDate"`
	DueDate                      string                        `xml:"cbc:DueDate"`
	InvoiceTypeCode              string                        `xml:"cbc:InvoiceTypeCode"`
	InvoicePeriod                *InvoicePeriod                `xml:"cac:InvoicePeriod,omitempty"`
	DocumentCurrencyCode         string                        `xml:"cbc:DocumentCurrencyCode"`
	OrderReference               *OrderReference               `xml:"cac:OrderReference,omitempty"`
	AdditionalDocumentReferences []AdditionalDocumentReference `xml:"cac:AdditionalDocumentReference,omitempty"`
	AccountingSupplierParty      AccountingSupplierParty       `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty      AccountingCustomerParty       `xml:"cac:AccountingCustomerParty"`
	PaymentMeans                 PaymentMeans                  `xml:"cac:PaymentMeans,omitempty"`
	TaxTotal                     TaxTotal                      `xml:"cac:TaxTotal"`
	LegalMonetaryTotal           LegalMonetaryTotal            `xml:"cac:LegalMonetaryTotal"`
	InvoiceLines                 []InvoiceLine                 `xml:"cac:InvoiceLine"`
}

type InvoicePeriod struct {
	StartDate string `xml:"cbc:StartDate"`
	EndDate   string `xml:"cbc:EndDate"`
}

type OrderReference struct {
	ID string `xml:"cbc:ID"`
}

type AccountingSupplierParty struct {
	Party Party `xml:"cac:Party"`
}

type AccountingCustomerParty struct {
	Party Party `xml:"cac:Party"`
}

type Party struct {
	EndpointID          *ID                   `xml:"cbc:EndpointID,omitempty"`
	PartyIdentification []PartyIdentification `xml:"cac:PartyIdentification,omitempty"`
	PartyName           *PartyName            `xml:"cac:PartyName,omitempty"`
	PostalAddress       Address               `xml:"cac:PostalAddress"`
	PartyTaxScheme      []PartyTaxScheme      `xml:"cac:PartyTaxScheme,omitempty"`
	PartyLegalEntity    PartyLegalEntity      `xml:"cac:PartyLegalEntity"`
	Contact             *Contact              `xml:"cac:Contact,omitempty"`
}

type ID struct {
	Value    string `xml:",chardata"`
	SchemeID string `xml:"schemeID,attr,omitempty"`
}

type PartyIdentification struct {
	ID ID `xml:"cbc:ID"`
}

type PartyName struct {
	Name string `xml:"cbc:Name"`
}

type Address struct {
	StreetName string  `xml:"cbc:StreetName"`
	CityName   string  `xml:"cbc:CityName"`
	PostalZone string  `xml:"cbc:PostalZone"`
	Country    Country `xml:"cac:Country"`
}

type Country struct {
	IdentificationCode string `xml:"cbc:IdentificationCode"`
}

type PartyTaxScheme struct {
	CompanyID string    `xml:"cbc:CompanyID"`
	TaxScheme TaxScheme `xml:"cac:TaxScheme"`
}

type TaxScheme struct {
	ID string `xml:"cbc:ID"`
}

type PartyLegalEntity struct {
	RegistrationName string `xml:"cbc:RegistrationName"`
	CompanyID        string `xml:"cbc:CompanyID"`
}

type Contact struct {
	Telephone      string `xml:"cbc:Telephone,omitempty"`
	ElectronicMail string `xml:"cbc:ElectronicMail,omitempty"`
}

type PaymentMeans struct {
	PaymentMeansCode      string           `xml:"cbc:PaymentMeansCode"`
	PayeeFinancialAccount FinancialAccount `xml:"cac:PayeeFinancialAccount"`
}

type FinancialAccount struct {
	ID string `xml:"cbc:ID"`
}

type TaxTotal struct {
	TaxAmount Money         `xml:"cbc:TaxAmount"`
	Subtotal  []TaxSubtotal `xml:"cac:TaxSubtotal"`
}

type TaxSubtotal struct {
	TaxableAmount Money       `xml:"cbc:TaxableAmount"`
	TaxAmount     Money       `xml:"cbc:TaxAmount"`
	TaxCategory   TaxCategory `xml:"cac:TaxCategory"`
}

type TaxCategory struct {
	ID        string    `xml:"cbc:ID"`
	Percent   string    `xml:"cbc:Percent"`
	TaxScheme TaxScheme `xml:"cac:TaxScheme"`
}

type Money struct {
	Value      string `xml:",chardata"`
	CurrencyID string `xml:"currencyID,attr"`
}

type LegalMonetaryTotal struct {
	LineExtensionAmount Money `xml:"cbc:LineExtensionAmount"`
	TaxExclusiveAmount  Money `xml:"cbc:TaxExclusiveAmount"`
	TaxInclusiveAmount  Money `xml:"cbc:TaxInclusiveAmount"`
	PayableAmount       Money `xml:"cbc:PayableAmount"`
}

type InvoiceLine struct {
	ID                  string   `xml:"cbc:ID"`
	InvoicedQuantity    Quantity `xml:"cbc:InvoicedQuantity"`
	LineExtensionAmount Money    `xml:"cbc:LineExtensionAmount"`
	Item                Item     `xml:"cac:Item"`
	Price               Price    `xml:"cac:Price"`
}

type Quantity struct {
	Value    string `xml:",chardata"`
	UnitCode string `xml:"unitCode,attr"`
}

type Item struct {
	Description           string                `xml:"cbc:Description"`
	Name                  string                `xml:"cbc:Name"`
	ClassifiedTaxCategory ClassifiedTaxCategory `xml:"cac:ClassifiedTaxCategory"`
}

type ClassifiedTaxCategory struct {
	ID        string    `xml:"cbc:ID"`
	Percent   string    `xml:"cbc:Percent"`
	TaxScheme TaxScheme `xml:"cac:TaxScheme"`
}

type Price struct {
	PriceAmount Money `xml:"cbc:PriceAmount"`
}

type AdditionalDocumentReference struct {
	ID               string        `xml:"cbc:ID"`
	DocumentTypeCode string        `xml:"cbc:DocumentTypeCode"`
	Attachment       UBLAttachment `xml:"cac:Attachment"`
}

type UBLAttachment struct {
	EmbeddedDocumentBinaryObject EmbeddedDocumentBinaryObject `xml:"cbc:EmbeddedDocumentBinaryObject"`
}

type EmbeddedDocumentBinaryObject struct {
	Value    string `xml:",chardata"`
	MimeCode string `xml:"mimeCode,attr"`
	Filename string `xml:"filename,attr"`
}

func cleanID(id string) string {
	var b strings.Builder
	for _, r := range id {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func Generate(inv *domain.Invoice, outputPath string) error {
	ublInv := Invoice{
		Xmlns:           "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2",
		XmlnsCac:        "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2",
		XmlnsCbc:        "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2",
		CustomizationID: "urn:cen.eu:en16931:2017#compliant#urn:fdc:peppol.eu:poacc:trns:invoice:3",
		ProfileID:       "urn:fdc:peppol.eu:poacc:bis:invoice:3",
		ID:              inv.Number,
		IssueDate:       inv.IssueDate.Format("2006-01-02"),
		InvoicePeriod: &InvoicePeriod{
			StartDate: inv.Period.Start.Format("2006-01-02"),
			EndDate:   inv.Period.End.Format("2006-01-02"),
		},
		DueDate:              inv.DueDate.Format("2006-01-02"),
		InvoiceTypeCode:      "380",
		DocumentCurrencyCode: inv.Currency,
		OrderReference: func() *OrderReference {
			if inv.OrderReference != "" {
				return &OrderReference{ID: inv.OrderReference}
			}
			return nil
		}(),
		AdditionalDocumentReferences: func() []AdditionalDocumentReference {
			var refs []AdditionalDocumentReference
			for _, att := range inv.Attachments {
				refs = append(refs, AdditionalDocumentReference{
					ID:               att.Filename,
					DocumentTypeCode: "916",
					Attachment: UBLAttachment{
						EmbeddedDocumentBinaryObject: EmbeddedDocumentBinaryObject{
							Value:    base64.StdEncoding.EncodeToString(att.Content),
							MimeCode: att.MimeType,
							Filename: att.Filename,
						},
					},
				})
			}
			return refs
		}(),
		AccountingSupplierParty: AccountingSupplierParty{
			Party: Party{
				EndpointID: &ID{Value: cleanID(inv.Supplier.CompanyID), SchemeID: "0208"}, // 0208 = Belgium CBE
				PartyIdentification: []PartyIdentification{
					{ID: ID{Value: inv.Supplier.CompanyID}},
				},
				PartyName: &PartyName{Name: inv.Supplier.Name},
				PostalAddress: Address{
					StreetName: inv.Supplier.Address.Street,
					CityName:   inv.Supplier.Address.City,
					PostalZone: inv.Supplier.Address.PostalCode,
					Country:    Country{IdentificationCode: inv.Supplier.Address.Country},
				},
				PartyTaxScheme: []PartyTaxScheme{
					{
						CompanyID: inv.Supplier.CompanyID,
						TaxScheme: TaxScheme{ID: "VAT"},
					},
				},
				PartyLegalEntity: PartyLegalEntity{
					RegistrationName: inv.Supplier.Name,
					CompanyID:        inv.Supplier.CompanyID,
				},
				Contact: &Contact{
					ElectronicMail: inv.Supplier.Email,
					Telephone:      inv.Supplier.Phone,
				},
			},
		},
		AccountingCustomerParty: AccountingCustomerParty{
			Party: Party{
				EndpointID: &ID{Value: cleanID(inv.Customer.CompanyID), SchemeID: "0208"},
				PartyIdentification: []PartyIdentification{
					{ID: ID{Value: inv.Customer.CompanyID}},
				},
				PartyName: &PartyName{Name: inv.Customer.Name},
				PostalAddress: Address{
					StreetName: inv.Customer.Address.Street,
					CityName:   inv.Customer.Address.City,
					PostalZone: inv.Customer.Address.PostalCode,
					Country:    Country{IdentificationCode: inv.Customer.Address.Country},
				},
				PartyTaxScheme: []PartyTaxScheme{
					{
						CompanyID: inv.Customer.CompanyID,
						TaxScheme: TaxScheme{ID: "VAT"},
					},
				},
				PartyLegalEntity: PartyLegalEntity{
					RegistrationName: inv.Customer.Name,
					CompanyID:        inv.Customer.CompanyID,
				},
			},
		},
		PaymentMeans: PaymentMeans{
			PaymentMeansCode: "30",
			PayeeFinancialAccount: FinancialAccount{
				ID: inv.Supplier.IBAN,
			},
		},
		TaxTotal: TaxTotal{
			TaxAmount: Money{Value: fmt.Sprintf("%.2f", inv.VATAmount), CurrencyID: inv.Currency},
			Subtotal: []TaxSubtotal{
				{
					TaxableAmount: Money{Value: fmt.Sprintf("%.2f", inv.Subtotal), CurrencyID: inv.Currency},
					TaxAmount:     Money{Value: fmt.Sprintf("%.2f", inv.VATAmount), CurrencyID: inv.Currency},
					TaxCategory: TaxCategory{
						ID:      "S",
						Percent: fmt.Sprintf("%.2f", inv.Lines[0].TaxPercent), // Assuming same tax for all for now
						TaxScheme: TaxScheme{
							ID: "VAT",
						},
					},
				},
			},
		},
		LegalMonetaryTotal: LegalMonetaryTotal{
			LineExtensionAmount: Money{Value: fmt.Sprintf("%.2f", inv.Subtotal), CurrencyID: inv.Currency},
			TaxExclusiveAmount:  Money{Value: fmt.Sprintf("%.2f", inv.Subtotal), CurrencyID: inv.Currency},
			TaxInclusiveAmount:  Money{Value: fmt.Sprintf("%.2f", inv.TotalAmount), CurrencyID: inv.Currency},
			PayableAmount:       Money{Value: fmt.Sprintf("%.2f", inv.TotalAmount), CurrencyID: inv.Currency},
		},
	}

	for i, line := range inv.Lines {
		ublInv.InvoiceLines = append(ublInv.InvoiceLines, InvoiceLine{
			ID:                  fmt.Sprintf("%d", i+1),
			InvoicedQuantity:    Quantity{Value: fmt.Sprintf("%.2f", line.Quantity), UnitCode: "DAY"},
			LineExtensionAmount: Money{Value: fmt.Sprintf("%.2f", line.NetAmount), CurrencyID: inv.Currency},
			Item: Item{
				Description: line.Description,
				Name:        line.ProjectCode,
				ClassifiedTaxCategory: ClassifiedTaxCategory{
					ID:      "S",
					Percent: fmt.Sprintf("%.2f", line.TaxPercent),
					TaxScheme: TaxScheme{
						ID: "VAT",
					},
				},
			},
			Price: Price{
				PriceAmount: Money{Value: fmt.Sprintf("%.2f", line.UnitPrice), CurrencyID: inv.Currency},
			},
		})
	}

	output, err := xml.MarshalIndent(ublInv, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal UBL XML: %w", err)
	}

	header := []byte(xml.Header)
	if err := os.WriteFile(outputPath, append(header, output...), 0644); err != nil {
		return fmt.Errorf("failed to write UBL file: %w", err)
	}

	return nil
}
