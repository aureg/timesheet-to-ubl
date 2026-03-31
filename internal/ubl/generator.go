package ubl

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"os"
	"ublcli/internal/domain"
)

type Invoice struct {
	XMLName                      xml.Name                      `xml:"urn:oasis:names:specification:ubl:schema:xsd:Invoice-2 Invoice"`
	CustomizationID              string                        `xml:"cbc:CustomizationID"`
	ProfileID                    string                        `xml:"cbc:ProfileID"`
	ID                           string                        `xml:"cbc:ID"`
	IssueDate                    string                        `xml:"cbc:IssueDate"`
	DueDate                      string                        `xml:"cbc:DueDate"`
	InvoiceTypeCode              string                        `xml:"cbc:InvoiceTypeCode"`
	DocumentCurrencyCode         string                        `xml:"cbc:DocumentCurrencyCode"`
	AccountingSupplierParty      AccountingSupplierParty       `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty      AccountingCustomerParty       `xml:"cac:AccountingCustomerParty"`
	PaymentMeans                 PaymentMeans                  `xml:"cac:PaymentMeans"`
	TaxTotal                     TaxTotal                      `xml:"cac:TaxTotal"`
	LegalMonetaryTotal           LegalMonetaryTotal            `xml:"cac:LegalMonetaryTotal"`
	InvoiceLines                 []InvoiceLine                 `xml:"cac:InvoiceLine"`
	AdditionalDocumentReferences []AdditionalDocumentReference `xml:"cac:AdditionalDocumentReference,omitempty"`
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
	PartyTaxScheme      *PartyTaxScheme       `xml:"cac:PartyTaxScheme,omitempty"`
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
	ElectronicMail string `xml:"cbc:ElectronicMail,omitempty"`
	Telephone      string `xml:"cbc:Telephone,omitempty"`
}

type PaymentMeans struct {
	PaymentMeansCode      string           `xml:"cbc:PaymentMeansCode"`
	PayeeFinancialAccount FinancialAccount `xml:"cac:PayeeFinancialAccount"`
}

type FinancialAccount struct {
	ID string `xml:"cbc:ID"`
}

type TaxTotal struct {
	TaxAmount Money `xml:"cbc:TaxAmount"`
	// Simplified TaxTotal for this exercise
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
	Description string `xml:"cbc:Description"`
	Name        string `xml:"cbc:Name"`
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

func Generate(inv *domain.Invoice, outputPath string) error {
	ublInv := Invoice{
		CustomizationID:      "urn:cen.eu:en16931:2017#compliant#urn:fdc:peppol.eu:poacc:trns:invoice:3",
		ProfileID:            "urn:fdc:peppol.eu:poacc:bis:invoice:3",
		ID:                   inv.Number,
		IssueDate:            inv.IssueDate.Format("2006-01-02"),
		DueDate:              inv.DueDate.Format("2006-01-02"),
		InvoiceTypeCode:      "380",
		DocumentCurrencyCode: inv.Currency,
		AccountingSupplierParty: AccountingSupplierParty{
			Party: Party{
				PartyName: &PartyName{Name: inv.Supplier.Name},
				PostalAddress: Address{
					StreetName: inv.Supplier.Address.Street,
					CityName:   inv.Supplier.Address.City,
					PostalZone: inv.Supplier.Address.PostalCode,
					Country:    Country{IdentificationCode: inv.Supplier.Address.Country},
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
				PartyName: &PartyName{Name: inv.Customer.Name},
				PostalAddress: Address{
					StreetName: inv.Customer.Address.Street,
					CityName:   inv.Customer.Address.City,
					PostalZone: inv.Customer.Address.PostalCode,
					Country:    Country{IdentificationCode: inv.Customer.Address.Country},
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
			},
			Price: Price{
				PriceAmount: Money{Value: fmt.Sprintf("%.2f", line.UnitPrice), CurrencyID: inv.Currency},
			},
		})
	}

	for _, att := range inv.Attachments {
		ublInv.AdditionalDocumentReferences = append(ublInv.AdditionalDocumentReferences, AdditionalDocumentReference{
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
