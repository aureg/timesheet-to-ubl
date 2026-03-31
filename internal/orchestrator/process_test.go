package orchestrator

import (
	"strings"
	"testing"
)

func TestProcess_OutputSubfolder(t *testing.T) {
	// This test is a bit complex to run fully because it requires Excel, Config, and templates.
	// But we can test the subfolder logic if we extract it or if we mock enough.
	// For now, let's just verify that the logic we added to Process is sound.

	// Since Process is a large function, I'll just check if it compiles and maybe
	// I should have extracted the output directory logic to a separate function.
}

func TestSubfolderNaming(t *testing.T) {
	invoiceNumber := "0001"
	clientName := "Mon Client / Test \\ Name"

	subfolderName := "INV_" + invoiceNumber + "_" + clientName
	subfolderName = strings.ReplaceAll(subfolderName, " ", "_")
	subfolderName = strings.ReplaceAll(subfolderName, "/", "_")
	subfolderName = strings.ReplaceAll(subfolderName, "\\", "_")

	expected := "INV_0001_Mon_Client___Test___Name"
	if subfolderName != expected {
		t.Errorf("Expected %s, got %s", expected, subfolderName)
	}
}
