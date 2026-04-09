package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ublcli/internal/config"
	"ublcli/internal/domain"
	"ublcli/internal/excel"
	"ublcli/internal/invoice"
	"ublcli/internal/pdf"
	excelrender "ublcli/internal/render/excel"
	"ublcli/internal/render/htmltmpl"
	"ublcli/internal/render/word"
	"ublcli/internal/ubl"
)

type GenerateOptions struct {
	ExcelPath         string
	SheetName         string
	ConfigPath        string
	TemplatePath      string
	ExcelTemplatePath string
	OutputDir         string
	InvoiceNumber     string
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
	inv, err := svc.Calculate(entries, cfg, opts.InvoiceNumber)
	if err != nil {
		return fmt.Errorf("calculating invoice: %w", err)
	}

	// 4. Resolve Template Path
	templatePath := config.ResolveTemplatePath(opts.TemplatePath, inv.InvoiceTemplate)
	if templatePath == "" {
		return fmt.Errorf("could not find invoice template (tried CLI option, config, and default locations)")
	}

	// 5. Determine and create output dir
	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = inv.OutputDir
	}
	if outputDir == "" {
		outputDir = "./dist"
	}

	// Append subfolder INV_{{InvoiceNumber}}_{{ClientName}}
	subfolderName := fmt.Sprintf("INV_%s_%s", inv.Number, inv.Customer.Name)
	// Sanitize subfolder name (remove characters that are invalid in file paths)
	subfolderName = strings.ReplaceAll(subfolderName, " ", "_")
	subfolderName = strings.ReplaceAll(subfolderName, "/", "_")
	subfolderName = strings.ReplaceAll(subfolderName, "\\", "_")

	outputDir = filepath.Join(outputDir, subfolderName)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	// 4. Render Invoice (HTML or Word)
	isWord := strings.ToLower(filepath.Ext(templatePath)) == ".docx"
	pdfPath := filepath.Join(outputDir, "invoice.pdf")
	absPdfPath, _ := filepath.Abs(pdfPath)

	if isWord {
		// Word flow
		docxPath := filepath.Join(outputDir, "invoice_filled.docx")
		wordRenderer := word.NewRenderer(templatePath)
		if err := wordRenderer.Render(inv, docxPath); err != nil {
			return fmt.Errorf("rendering word: %w", err)
		}

		// Convert docx to pdf using MS Office
		renderer := pdf.NewMsOfficeRenderer()
		absDocxPath, _ := filepath.Abs(docxPath)
		if err := renderer.Convert(absDocxPath, outputDir); err != nil {
			fmt.Printf("Warning: Word to PDF conversion failed (MS Office): %v\n", err)
		} else {
			// Find the produced PDF
			producedPath := filepath.Join(outputDir, "invoice_filled.pdf")
			// Rename it to invoice.pdf
			if err := os.Rename(producedPath, pdfPath); err != nil {
				fmt.Printf("Warning: Failed to rename word-produced PDF: %v\n", err)
			}
		}
	} else {
		// HTML flow
		htmlPath := filepath.Join(outputDir, "invoice.html")
		renderer := htmltmpl.NewRenderer(templatePath)
		if err := renderer.Render(inv, htmlPath); err != nil {
			return fmt.Errorf("rendering html: %w", err)
		}

		// Convert HTML to PDF using Chrome
		pdfRenderer := pdf.NewChromeRenderer()
		absHtmlPath, _ := filepath.Abs(htmlPath)
		if err := pdfRenderer.Convert(absHtmlPath, absPdfPath); err != nil {
			fmt.Printf("Warning: Invoice PDF conversion failed: %v\n", err)
		}
	}

	// 5. Render Excel Template if provided
	if opts.ExcelTemplatePath != "" {
		excelOutPath := filepath.Join(outputDir, "timesheet_filled.xlsx")
		if strings.HasSuffix(strings.ToLower(opts.ExcelTemplatePath), ".xlsm") {
			excelOutPath = filepath.Join(outputDir, "timesheet_filled.xlsm")
		}
		excelRenderer := excelrender.NewRenderer(opts.ExcelTemplatePath)
		if err := excelRenderer.Render(inv, entries, excelOutPath); err != nil {
			fmt.Printf("Warning: Excel template rendering failed: %v\n", err)
		} else {
			// Attach filled excel
			if data, err := os.ReadFile(excelOutPath); err == nil {
				inv.Attachments = append(inv.Attachments, domain.Attachment{
					Filename: filepath.Base(excelOutPath),
					MimeType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
					Content:  data,
				})
			}
		}
	}

	// 6. Attach Invoice PDF
	if data, err := os.ReadFile(pdfPath); err == nil {
		inv.Attachments = append(inv.Attachments, domain.Attachment{
			Filename: "invoice.pdf",
			MimeType: "application/pdf",
			Content:  data,
		})
	}

	// 6b. Attach original Excel as PDF using MS Office
	absExcelPath, _ := filepath.Abs(opts.ExcelPath)
	if _, err := os.Stat(absExcelPath); err != nil {
		fmt.Printf("Warning: Excel file not found for PDF conversion: %v\n", err)
	} else {
		msoRenderer := pdf.NewMsOfficeRenderer()
		if err := msoRenderer.Convert(absExcelPath, outputDir); err != nil {
			fmt.Printf("Warning: Excel to PDF conversion failed (MS Office): %v\n", err)
		} else {
			// MS Office outputs <basename>.pdf in the output dir
			base := filepath.Base(absExcelPath)
			producedName := strings.TrimSuffix(base, filepath.Ext(base)) + ".pdf"
			producedPath := filepath.Join(outputDir, producedName)
			if data, err := os.ReadFile(producedPath); err == nil {
				inv.Attachments = append(inv.Attachments, domain.Attachment{
					Filename: "timesheet.pdf",
					MimeType: "application/pdf",
					Content:  data,
				})
			} else {
				fmt.Printf("Warning: Could not read produced timesheet PDF: %v\n", err)
			}
		}
	}

	// 7. Generate UBL
	ublPath := filepath.Join(outputDir, "invoice-ubl.xml")
	if err := ubl.Generate(inv, ublPath); err != nil {
		return fmt.Errorf("generating ubl: %w", err)
	}

	return nil
}
