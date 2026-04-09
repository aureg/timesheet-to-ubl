# Handoff Document - ublcli

This document contains instructions and information necessary to continue the development of the `ublcli` project on
another computer or with another IA.

## Project Overview

`ublcli` is a Go-based tool (v1.26+) designed to generate invoices compliant with the UBL XML standard, along with HTML,
Word, and PDF renders, starting from a daily activity report (timesheet) in Excel format.

### Main Features:

- **Excel Import**: Flexible reading (FR/EN/Numeric dates, numbers with commas or dots).
- **Business Logic**: Aggregation of services by project, calculation of Net/VAT/Total amounts.
- **Config Enrichment**: Configuration resolution system by priorities (Client+Project > Client > Project > Default)
  loaded via YAML.
- **Multi-format Rendering**: Generation of HTML (via `html/template`) or Word (via `{{Key}}` placeholders), and UBL
  XML.
- **Excel Template Support**: Filling a specific Excel template (e.g., for timesheet validation) and converting it to
  PDF.
- **PDF Conversion**:
    - HTML -> PDF via Chrome/Chromium.
    - Word/Excel -> PDF via MS Office (PowerShell COM Interop).
- **Dual Interface**: Standard CLI mode and interactive TUI (Bubble Tea).

## Technical Architecture

The project follows a modular and testable architecture:

- `cmd/ublcli`: Main entry point.
- `internal/domain`: Business data structures (`Invoice`, `Party`, `Attachment`, etc.).
- `internal/config`: Loading and resolution of YAML configuration.
- `internal/excel`: Excel importer using `excelize`.
- `internal/invoice`: Service for amount calculation and aggregation.
- `internal/render/htmltmpl`: HTML invoice generator.
- `internal/render/word`: Word template (.docx) invoice generator.
- `internal/render/excel`: Excel template (.xlsx, .xlsm) filler.
- `internal/pdf`: Wrappers for PDF conversion (Chrome and MS Office).
- `internal/ubl`: UBL 2.1 XML generator.
- `internal/orchestrator`: Coordination of the full generation workflow.
- `internal/cli`: Command and flag management.
- `internal/tui`: Interactive user interface.

## Development Environment

### Default Configuration and Resolution:

The application automatically looks for its configuration file in:
`%userHome%\.timesheet2ubl\config.yaml` (Windows) or `~/.timesheet2ubl\config.yaml` (Linux/macOS).

**Template Resolution:**
The tool resolves the invoice template path in this order of priority:

1. `--template` CLI option.
2. `invoice_template` value for the client in the YAML config.
3. `invoice_template` value in the `default` section of the YAML config.
4. `invoice.html` or `invoice.docx` file in `%userHome%\.timesheet2ubl\`.

**Timesheet Template Resolution:**
Similarly, for the optional Excel timesheet template:

1. `--excel-template` CLI option.
2. `excel_template` value for the client in the YAML config.
3. `excel_template` value in the `default` section.
4. `.timesheet2ubl/` folder search.

### Prerequisites:

- **Go 1.26** or higher.
- **Microsoft Office** (Excel/Word): Necessary for Word/Excel to PDF conversion via PowerShell COM.
- **Google Chrome** or **Chromium**: For HTML to PDF conversion.
- A Go-compatible editor (JetBrains GoLand recommended).

### Major Dependencies:

- `github.com/nguyenthenguyen/docx`: Word file manipulation.
- `github.com/xuri/excelize/v2`: Excel file manipulation.
- `gopkg.in/yaml.v3`: YAML parsing.
- `github.com/charmbracelet/bubbletea`: TUI framework.

## Installation and Usage

1. **Initialization**:
   ```powershell
   go mod download
   ```
2. **Compilation**:
   ```powershell
   go build -o ublcli.exe ./cmd/ublcli/main.go
   ```
3. **CLI Execution**:
   ```powershell
   # Basic usage with default config and template
   .\ublcli.exe generate --timesheet-in data.xlsx --invoice-number 0001
   
   # Using aliases
   .\ublcli.exe generate --tsin data.xlsx --inv-num 0001
   
   # Specifying a template
   .\ublcli.exe generate --tsin data.xlsx --inv-num 0001 --template my-invoice.docx
   ```
4. **TUI Execution**:
   ```powershell
   .\ublcli.exe tui
   ```

## Template Documentation

### Common Placeholders (HTML & Word)

| Placeholder         | Description             | Format / Example |
|:--------------------|:------------------------|:-----------------|
| `{{InvoiceNumber}}` | Full invoice number     | `2026-0001`      |
| `{{InvoiceDate}}`   | Issuance date           | `31/03/2026`     |
| `{{DueDate}}`       | Due date                | `30/04/2026`     |
| `{{PeriodStart}}`   | Start of service period | `01/03/2026`     |
| `{{PeriodEnd}}`     | End of service period   | `31/03/2026`     |
| `{{ClientName}}`    | Client name             | `Client Name`    |
| `{{TotalHours}}`    | Total hours sum         | `145.50`         |
| `{{Subtotal}}`      | Total Net amount        | `9457.50`        |
| `{{VATAmount}}`     | Total VAT amount        | `1986.08`        |
| `{{TotalAmount}}`   | Total amount Incl. VAT  | `11443.58`       |
| `{{DailyRate}}`     | Daily rate              | `650.00`         |
| `{{ManagerName}}`   | Manager Name            | `John Doe`       |

### Excel Template Specifics (`internal/render/excel`)

When an Excel template is used (via `--excel-template`), the following cells are filled:

- **I8**: Consultant Name
- **I9**: Supplier Name
- **I11**: First day of the billing month
- **D57**: Manager Name (from config)
- **R10**: Daily Rate
- **Column D (from row 13)**: Hours per day (D13 = Day 1, D14 = Day 2, etc.)

The resulting file is converted to PDF in Landscape mode and attached to the UBL.

## Key Points for Future Development

1. **Daily Rate**: The parameter is now `daily_rate` in YAML and `DailyRate` in the code (replacing the old HourlyRate).
2. **PDF Conversion**: The MS Office renderer uses a PowerShell script to control Excel/Word via COM. It forces
   Landscape orientation for Excel.
3. **UBL Attachments**: The `invoice.pdf` is attached first, followed by `timesheet.pdf` (if generated).
4. **Traceability**: The source Excel file is always copied to the output directory as `source_<filename>`.
5. **Language**: All comments, logs, and TUI labels must remain in English.

## Future Improvements

- Multi-currency support per line.
- Logo integration in the default HTML template.
- Direct Peppol/Platform upload.
- Better error handling when Chrome or Office is missing.
