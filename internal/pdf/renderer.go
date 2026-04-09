package pdf

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Renderer interface {
	Convert(htmlPath, pdfPath string) error
}

type ChromeRenderer struct {
	Timeout time.Duration
}

type LibreOfficeRenderer struct {
	Timeout time.Duration
}

func NewChromeRenderer() *ChromeRenderer {
	return &ChromeRenderer{
		Timeout: 30 * time.Second,
	}
}

func NewLibreOfficeRenderer() *LibreOfficeRenderer {
	return &LibreOfficeRenderer{
		Timeout: 60 * time.Second,
	}
}

func (r *ChromeRenderer) Convert(htmlPath, pdfPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	// Try to find chrome or chromium
	paths := []string{"chrome", "google-chrome", "chromium", "chromium-browser", `C:\Program Files\Google\Chrome\Application\chrome.exe`}

	var cmdPath string
	for _, p := range paths {
		if path, err := exec.LookPath(p); err == nil {
			cmdPath = path
			break
		}
	}

	if cmdPath == "" {
		// Fallback to hardcoded windows path if LookPath failed but file exists
		cmdPath = `C:\Program Files\Google\Chrome\Application\chrome.exe`
	}

	args := []string{
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--no-pdf-header-footer",
		"--run-all-compositor-stages-before-draw",
		"--virtual-time-budget=10000",
		"--hide-scrollbars",
		"--disable-breakpad",
		"--disable-extensions",
		"--disable-infobars",
		"--disable-dev-shm-usage",
		"--window-size=1200,1600",
		fmt.Sprintf("--print-to-pdf=%s", pdfPath),
		htmlPath,
	}

	cmd := exec.CommandContext(ctx, cmdPath, args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("chrome conversion failed: %w (output: %s)", err, string(output))
	}

	return nil
}

type MsOfficeRenderer struct {
	Timeout time.Duration
}

func NewMsOfficeRenderer() *MsOfficeRenderer {
	return &MsOfficeRenderer{
		Timeout: 60 * time.Second,
	}
}

func (r *MsOfficeRenderer) Convert(inputPath, outputDir string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	absInputPath, _ := filepath.Abs(inputPath)
	ext := strings.ToLower(filepath.Ext(absInputPath))
	base := filepath.Base(absInputPath)
	producedName := strings.TrimSuffix(base, filepath.Ext(base)) + ".pdf"
	absOutputPath := filepath.Join(outputDir, producedName)
	absOutputPath, _ = filepath.Abs(absOutputPath)

	var psScript string
	if ext == ".docx" || ext == ".doc" {
		// Word conversion script
		psScript = fmt.Sprintf(`
$word = New-Object -ComObject Word.Application
$word.Visible = $false
try {
    $doc = $word.Documents.Open("%s")
    $doc.SaveAs([ref]"%s", [ref]17) # 17 is wdExportFormatPDF
    $doc.Close([ref]0) # 0 is wdDoNotSaveChanges
} finally {
    $word.Quit()
    [System.Runtime.Interopservices.Marshal]::ReleaseComObject($word) | Out-Null
    [GC]::Collect()
    [GC]::WaitForPendingFinalizers()
}
`, absInputPath, absOutputPath)
	} else if ext == ".xlsx" || ext == ".xls" {
		// Excel conversion script
		psScript = fmt.Sprintf(`
$excel = New-Object -ComObject Excel.Application
$excel.Visible = $false
$excel.DisplayAlerts = $false
try {
    $wb = $excel.Workbooks.Open("%s")
    $wb.ExportAsFixedFormat(0, "%s") # 0 is xlTypePDF
    $wb.Close($false)
} finally {
    $excel.Quit()
    [System.Runtime.Interopservices.Marshal]::ReleaseComObject($excel) | Out-Null
    [GC]::Collect()
    [GC]::WaitForPendingFinalizers()
}
`, absInputPath, absOutputPath)
	} else {
		return fmt.Errorf("unsupported file extension for MS Office conversion: %s", ext)
	}

	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ms office conversion failed: %w (output: %s)", err, string(output))
	}

	return nil
}

func (r *LibreOfficeRenderer) Convert(inputPath, outputDir string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	// Path to soffice.exe
	cmdPath := `C:\Program Files\LibreOffice\program\soffice.exe`

	args := []string{
		"--headless",
		"--convert-to", "pdf",
		"--outdir", outputDir,
		inputPath,
	}

	cmd := exec.CommandContext(ctx, cmdPath, args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("libreoffice conversion failed: %w (output: %s)", err, string(output))
	}

	return nil
}
