package ubl

import (
	"os"
	"testing"
	"time"
	"ublcli/internal/domain"
)

func TestGenerate(t *testing.T) {
	inv := &domain.Invoice{
		Number:    "INV-2023001",
		IssueDate: time.Now(),
		DueDate:   time.Now().AddDate(0, 0, 30),
		Currency:  "EUR",
		Supplier: domain.Party{
			Name:      "Supplier Co",
			CompanyID: "BE0123456789",
			Address: domain.Address{
				Street:     "Street 1",
				City:       "Brussels",
				PostalCode: "1000",
				Country:    "BE",
			},
		},
		Customer: domain.Party{
			Name:      "Customer Co",
			CompanyID: "BE0987654321",
			Address: domain.Address{
				Street:     "Avenue 2",
				City:       "Namur",
				PostalCode: "5000",
				Country:    "BE",
			},
		},
		Lines: []domain.InvoiceLine{
			{
				ProjectCode: "P1",
				Description: "Consulting",
				Quantity:    1.0,
				UnitPrice:   1000.0,
				TaxPercent:  21.0,
				NetAmount:   1000.0,
				TaxAmount:   210.0,
				GrossAmount: 1210.0,
			},
		},
		Subtotal:                1000.0,
		VATAmount:               210.0,
		TotalAmount:             1210.0,
		OrderReference:          "PO-456",
		StructuredCommunication: "+++2023/0001/0013+++",
		Attachments: []domain.Attachment{
			{Filename: "test.pdf", MimeType: "application/pdf", Content: []byte("fake")},
			{Filename: "timesheet.pdf", MimeType: "application/pdf", Content: []byte("fake2")},
		},
	}

	tmpFile := "test-ubl.xml"
	defer os.Remove(tmpFile)

	err := Generate(inv, tmpFile)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	xmlStr := string(content)

	// Check if both attachments are present
	if !contains(xmlStr, `<cbc:ID>test.pdf</cbc:ID>`) {
		t.Error("XML missing first attachment: test.pdf")
	}
	if !contains(xmlStr, `<cbc:ID>timesheet.pdf</cbc:ID>`) {
		t.Error("XML missing second attachment: timesheet.pdf")
	}

	// Check namespaces
	expectedNS := []string{
		`xmlns="urn:oasis:names:specification:ubl:schema:xsd:Invoice-2"`,
		`xmlns:cac="urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2"`,
		`xmlns:cbc="urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2"`,
	}

	for _, ns := range expectedNS {
		if !contains(xmlStr, ns) {
			t.Errorf("XML missing namespace: %s", ns)
		}
	}

	// Check Peppol IDs
	if !contains(xmlStr, `<cbc:CustomizationID>urn:cen.eu:en16931:2017#compliant#urn:fdc:peppol.eu:poacc:trns:invoice:3</cbc:CustomizationID>`) {
		t.Error("XML missing correct CustomizationID")
	}

	// Check TaxSubtotal
	if !contains(xmlStr, `<cac:TaxSubtotal>`) {
		t.Error("XML missing TaxSubtotal")
	}

	// Check Order of elements (simplified check)
	idxDocRef := index(xmlStr, "<cac:AdditionalDocumentReference>")
	idxSupplier := index(xmlStr, "<cac:AccountingSupplierParty>")

	// DocRef must be before AccountingSupplierParty
	if idxDocRef != -1 && idxSupplier != -1 && idxDocRef > idxSupplier {
		t.Error("cac:AdditionalDocumentReference must appear before cac:AccountingSupplierParty")
	}

	// Check EndpointID formatting
	if !contains(xmlStr, `<cbc:EndpointID schemeID="0208">0123456789</cbc:EndpointID>`) {
		t.Error("Supplier EndpointID not correctly formatted (should be only digits)")
	}
	if !contains(xmlStr, `<cbc:EndpointID schemeID="0208">0987654321</cbc:EndpointID>`) {
		t.Error("Customer EndpointID not correctly formatted (should be only digits)")
	}

	// Check OrderReference
	if !contains(xmlStr, "<cac:OrderReference>") || !contains(xmlStr, "<cbc:ID>PO-456</cbc:ID>") {
		t.Errorf("XML missing OrderReference. Got: %s", xmlStr)
	}

	// Check StructuredCommunication (PaymentID)
	if !contains(xmlStr, "<cbc:PaymentID>+++2023/0001/0013+++</cbc:PaymentID>") {
		t.Errorf("XML missing PaymentID (StructuredCommunication). Got: %s", xmlStr)
	}
}

func index(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func contains(s, substr string) bool {
	return index(s, substr) != -1
}
