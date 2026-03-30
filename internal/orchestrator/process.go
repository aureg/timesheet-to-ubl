package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"

	"ublcli/internal/config"
	"ublcli/internal/domain"
	"ublcli/internal/excel"
	"ublcli/internal/invoice"
	"ublcli/internal/pdf"
	"ublcli/internal/render/htmltmpl"
	"ublcli/internal/ubl"
)

type GenerateOptions struct {
	ExcelPath    string
	SheetName    string
	ConfigPath   string
	TemplatePath string
	OutputDir    string
}

func Process(opts GenerateOptions) error {
	// 1. Load config
	cfg, err := config.LoadConfig(opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// 2. Import Excel
	imp := excel.NewImporter()
	entries, err := imp.Import(opts.ExcelPath, opts.SheetName)
	if err != nil {
		return fmt.Errorf("importing excel: %w", err)
	}

	// 3. Calculate Invoice
	svc := invoice.NewService()
	inv, err := svc.Calculate(entries, cfg)
	if err != nil {
		return fmt.Errorf("calculating invoice: %w", err)
	}

	// Create output dir
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	// 4. Render HTML
	htmlPath := filepath.Join(opts.OutputDir, "invoice.html")
	renderer := htmltmpl.NewRenderer(opts.TemplatePath)
	if err := renderer.Render(inv, htmlPath); err != nil {
		return fmt.Errorf("rendering html: %w", err)
	}

	// 5. Convert PDF
	pdfPath := filepath.Join(opts.OutputDir, "invoice.pdf")
	pdfRenderer := pdf.NewChromeRenderer()

	absHtmlPath, _ := filepath.Abs(htmlPath)
	absPdfPath, _ := filepath.Abs(pdfPath)

	if err := pdfRenderer.Convert(absHtmlPath, absPdfPath); err != nil {
		fmt.Printf("Warning: PDF conversion failed: %v\n", err)
	} else {
		// 6. Attach PDF to Invoice
		pdfContent, err := os.ReadFile(pdfPath)
		if err == nil {
			inv.Attachments = append(inv.Attachments, domain.Attachment{
				Filename: "invoice.pdf",
				MimeType: "application/pdf",
				Content:  pdfContent,
			})
		}
	}

	// 7. Generate UBL
	ublPath := filepath.Join(opts.OutputDir, "invoice-ubl.xml")
	if err := ubl.Generate(inv, ublPath); err != nil {
		return fmt.Errorf("generating ubl: %w", err)
	}

	return nil
}
