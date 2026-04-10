package cli

import (
	"flag"
	"fmt"
	"os"

	"ublcli/internal/config"
	"ublcli/internal/orchestrator"
	"ublcli/internal/tui"
)

func Run() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "generate":
		handleGenerate()
	case "tui":
		if err := tui.Start(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: ublcli <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  generate  Generate invoice files from Excel")
	fmt.Println("  tui       Launch terminal user interface")
}

func handleGenerate() {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	opts := orchestrator.GenerateOptions{}

	fs.StringVar(&opts.TimesheetInPath, "timesheet-in", "", "Path to timesheet Excel file")
	fs.StringVar(&opts.TimesheetInPath, "tsin", "", "Alias for --timesheet-in")
	fs.StringVar(&opts.SheetName, "sheet", "Data", "Excel sheet name")
	defaultConfigPath, _ := config.GetDefaultConfigPath()
	fs.StringVar(&opts.ConfigPath, "config", defaultConfigPath, "Path to YAML config file")
	fs.StringVar(&opts.TemplatePath, "template", "", "Path to HTML or Word (.docx) template (or name in .timesheet2ubl/)")
	fs.StringVar(&opts.ExcelTemplatePath, "excel-template", "", "Path to Excel (.xlsx, .xlsm) template")
	fs.StringVar(&opts.OutputDir, "out", "", "Output directory (default: ./dist or from config)")
	fs.StringVar(&opts.InvoiceNum, "invoice-number", "", "Invoice number (4 digits, e.g., 0001)")
	fs.StringVar(&opts.InvoiceNum, "inv-num", "", "Alias for --invoice-number")

	fs.Parse(os.Args[2:])

	if opts.TimesheetInPath == "" || opts.ConfigPath == "" || opts.InvoiceNum == "" {
		fmt.Println("Error: --timesheet-in, --config, and --invoice-number are required")
		fs.Usage()
		os.Exit(1)
	}

	if len(opts.InvoiceNum) != 4 {
		fmt.Println("Error: --invoice-number must be exactly 4 digits")
		os.Exit(1)
	}

	if err := orchestrator.Process(opts); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generation successful!")
}
