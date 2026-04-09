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

	// 1. Get first day of billing month
	// Based on the instructions: "le premier jour du mois de facturation"
	// We use the start of the invoice period.
	periodStart := invoice.Period.Start
	firstDay := time.Date(periodStart.Year(), periodStart.Month(), 1, 0, 0, 0, 0, periodStart.Location())

	// 2. Fill cell I11
	// The instructions don't specify the sheet, so we'll try the first one or a default
	// If the template is Template_TimeSheet_NRB_V9.1.xlsm, it might have a specific sheet name.
	// But usually, it's the active sheet or the first one.
	sheet := f.GetSheetName(0)
	if sheet == "" {
		sheet = "Sheet1"
	}

	// We set the value as time.Time so Excel can apply its own cell formatting (e.g., "mmmm yyyy")
	if err := f.SetCellValue(sheet, "I11", firstDay); err != nil {
		return fmt.Errorf("failed to set cell I11: %w", err)
	}

	// 2.5 Fill Hours in column D from line 13
	// D13 = day 1, D14 = day 2, ...
	// Aggregate hours by day
	hoursByDay := make(map[int]float64)
	for _, entry := range entries {
		if entry.Date.Month() == periodStart.Month() && entry.Date.Year() == periodStart.Year() {
			hoursByDay[entry.Date.Day()] += entry.Hours
		}
	}

	// Find number of days in the month
	lastDayOfMonth := time.Date(periodStart.Year(), periodStart.Month()+1, 0, 0, 0, 0, 0, periodStart.Location()).Day()

	for day := 1; day <= lastDayOfMonth; day++ {
		row := 13 + day - 1
		cell := fmt.Sprintf("D%d", row)
		hours := hoursByDay[day]

		if hours > 0 {
			// Format as d,dd (French format with comma)
			// We can set it as a string or try to set a number format.
			// The user said "affiche le total des heures au format d,dd".
			// If we want Excel to treat it as a number but display with comma, we should use a style.
			// But for simplicity and based on "d,dd", a string might be what they expect if the template isn't pre-formatted.
			// However, excelize can set float and we can set a custom number format.

			// Let's try setting the value as float and applying a format if possible,
			// or just format it as a string with comma.
			if err := f.SetCellValue(sheet, cell, hours); err != nil {
				return fmt.Errorf("failed to set cell %s: %w", cell, err)
			}
		}
	}

	// 2.7 Force recalculation on load
	// Some formulas might not be calculated when opening the file.
	// We tell Excel to recalculate everything.
	// In excelize v2.10.1, we use SetCalcProps.
	// We need a pointer to bool, and since BoolPtr is internal, we create one.
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
