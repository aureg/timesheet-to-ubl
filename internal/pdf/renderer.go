package pdf

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type Renderer interface {
	Convert(htmlPath, pdfPath string) error
}

type ChromeRenderer struct {
	Timeout time.Duration
}

func NewChromeRenderer() *ChromeRenderer {
	return &ChromeRenderer{
		Timeout: 30 * time.Second,
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
		fmt.Sprintf("--print-to-pdf=%s", pdfPath),
		htmlPath,
	}

	cmd := exec.CommandContext(ctx, cmdPath, args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("chrome conversion failed: %w (output: %s)", err, string(output))
	}

	return nil
}
