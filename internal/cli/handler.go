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

	fs.StringVar(&opts.ExcelPath, "excel", "", "Path to Excel file")
	fs.StringVar(&opts.SheetName, "sheet", "Data", "Excel sheet name")
	defaultConfigPath, _ := config.GetDefaultConfigPath()
	fs.StringVar(&opts.ConfigPath, "config", defaultConfigPath, "Path to YAML config file")
	fs.StringVar(&opts.TemplatePath, "template", "", "Path to HTML or Word (.docx) template (or name in .timesheet2ubl/)")
	fs.StringVar(&opts.OutputDir, "out", "", "Output directory (default: ./dist or from config)")
	fs.StringVar(&opts.InvoiceNumber, "number", "", "Invoice number (4 digits, e.g., 0001)")

	fs.Parse(os.Args[2:])

	if opts.ExcelPath == "" || opts.ConfigPath == "" || opts.InvoiceNumber == "" {
		fmt.Println("Error: --excel, --config, and --number are required")
		fs.Usage()
		os.Exit(1)
	}

	if len(opts.InvoiceNumber) != 4 {
		fmt.Println("Error: --number must be exactly 4 digits")
		os.Exit(1)
	}

	if err := orchestrator.Process(opts); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generation successful!")
}
