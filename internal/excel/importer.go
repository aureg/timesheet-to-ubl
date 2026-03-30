package excel

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"ublcli/internal/domain"
)

type Importer struct{}

func NewImporter() *Importer {
	return &Importer{}
}

func (i *Importer) Import(filepath, sheetName string) ([]domain.TimesheetEntry, error) {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet %s: %w", sheetName, err)
	}

	if len(rows) < 1 {
		return nil, fmt.Errorf("sheet is empty")
	}

	header := rows[0]
	colMap := make(map[string]int)
	for idx, col := range header {
		colMap[col] = idx
	}

	requiredCols := []string{"Date", "Projet", "Désignation du projet", "Tâche", "Désignation de la tâche", "H.", "Sp.", "Ticket", "Commentaire"}
	for _, rc := range requiredCols {
		if _, ok := colMap[rc]; !ok {
			return nil, fmt.Errorf("missing required column: %s", rc)
		}
	}

	var entries []domain.TimesheetEntry
	for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]
		if len(row) == 0 {
			continue
		}

		entry := domain.TimesheetEntry{
			LineNumber: rowIdx + 1,
		}

		// Helper to get col value safely
		getVal := func(name string) string {
			idx := colMap[name]
			if idx < len(row) {
				return row[idx]
			}
			return ""
		}

		dateStr := getVal("Date")
		if dateStr == "" {
			continue // skip empty rows
		}

		// Excelize might return date as serial number or formatted string
		t, err := i.parseDate(dateStr)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid date '%s': %w", entry.LineNumber, dateStr, err)
		}
		entry.Date = t

		entry.Project = getVal("Projet")
		entry.ProjectDesc = getVal("Désignation du projet")
		entry.Task = getVal("Tâche")
		entry.TaskDesc = getVal("Désignation de la tâche")
		if idx, ok := colMap["Nom du client"]; ok && idx < len(row) {
			entry.ClientName = row[idx]
		}
		entry.Ticket = getVal("Ticket")
		entry.Comment = getVal("Commentaire")

		hStr := getVal("H.")
		if hStr != "" {
			// Replace comma with dot for float parsing (common in French Excel)
			hStr = strings.ReplaceAll(hStr, ",", ".")
			h, err := strconv.ParseFloat(hStr, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid hours '%s': %w", entry.LineNumber, hStr, err)
			}
			entry.Hours = h
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (i *Importer) parseDate(val string) (time.Time, error) {
	// Try standard formats
	formats := []string{
		"2006-01-02",
		"02/01/2006",
		"02-01-2006",
		"1/2/06",
		"01/02/06",
		"02/01/06",
		"01-02-06",
		"02-01-06",
	}
	for _, f := range formats {
		t, err := time.Parse(f, val)
		if err == nil {
			return t, nil
		}
	}

	// Try excel serial
	f, err := strconv.ParseFloat(val, 64)
	if err == nil {
		return excelize.ExcelDateToTime(f, false)
	}

	return time.Time{}, fmt.Errorf("unknown date format")
}
