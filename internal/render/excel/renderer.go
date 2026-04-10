package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ublcli/internal/domain"

	"github.com/xuri/excelize/v2"
)

type Renderer struct {
	templatePath string
}

func NewRenderer(templatePath string) *Renderer {
	return &Renderer{
		templatePath: templatePath,
	}
}

func (r *Renderer) Render(invoice *domain.Invoice, entries []domain.TimesheetEntry, outputPath string) error {
	absPath, _ := filepath.Abs(r.templatePath)
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to access excel template '%s' (absolute path: %s): %w", r.templatePath, absPath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("excel template path is a directory, not a file: '%s' (absolute path: %s)", r.templatePath, absPath)
	}

	f, err := excelize.OpenFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to open excel template: %w", err)
	}
	defer f.Close()

	// 1. First day of billing month
	periodStart := invoice.Period.Start
	firstDay := time.Date(periodStart.Year(), periodStart.Month(), 1, 0, 0, 0, 0, periodStart.Location())

	// 2. Fill header cells
	// Get the first sheet of the workbook
	sheet := f.GetSheetName(0)
	if sheet == "" {
		sheet = "Sheet1"
	}

	// I11: first day of month (date)
	if err := f.SetCellValue(sheet, "I11", firstDay); err != nil {
		return fmt.Errorf("failed to set cell I11: %w", err)
	}

	// I8: Consultant Name
	if err := f.SetCellValue(sheet, "I8", invoice.ConsultantName); err != nil {
		return fmt.Errorf("failed to set cell I8: %w", err)
	}

	// I9: Supplier Name
	if err := f.SetCellValue(sheet, "I9", invoice.Supplier.Name); err != nil {
		return fmt.Errorf("failed to set cell I9: %w", err)
	}

	// D57: Manager Name
	if err := f.SetCellValue(sheet, "D57", invoice.ManagerName); err != nil {
		return fmt.Errorf("failed to set cell D57: %w", err)
	}

	// R10 : Daily Rate (based on the unit price of the first line)
	dailyRate := 0.0
	if len(invoice.Lines) > 0 {
		dailyRate = invoice.Lines[0].UnitPrice
	}
	if err := f.SetCellValue(sheet, "R10", dailyRate); err != nil {
		return fmt.Errorf("failed to set cell R10: %w", err)
	}

	// 2.5 Fill daily hours (Column D, starting from row 13)
	// Aggregate hours by day for the current month
	hoursByDay := make(map[int]float64)
	for _, entry := range entries {
		if entry.Date.Month() == periodStart.Month() && entry.Date.Year() == periodStart.Year() {
			hoursByDay[entry.Date.Day()] += entry.Hours
		}
	}

	// Number of days in the month
	lastDayOfMonth := time.Date(periodStart.Year(), periodStart.Month()+1, 0, 0, 0, 0, 0, periodStart.Location()).Day()

	for day := 1; day <= lastDayOfMonth; day++ {
		row := 13 + day - 1
		cell := fmt.Sprintf("D%d", row)
		hours := hoursByDay[day]

		if hours > 0 {
			// Write the raw numeric value, Excel will handle formatting via the template
			if err := f.SetCellValue(sheet, cell, hours); err != nil {
				return fmt.Errorf("failed to set cell %s: %w", cell, err)
			}
		}
	}

	// 2.7 Force formula recalculation on file open
	forceRecalc := true
	if err := f.SetCalcProps(&excelize.CalcPropsOptions{
		FullCalcOnLoad: &forceRecalc,
		ForceFullCalc:  &forceRecalc,
	}); err != nil {
		return fmt.Errorf("failed to set CalcProps: %w", err)
	}

	// 3. Save to output path
	if err := f.SaveAs(outputPath); err != nil {
		return fmt.Errorf("failed to save excel: %w", err)
	}

	return nil
}
