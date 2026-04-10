# timesheet-to-ubl (ublcli)

`ublcli` is a powerful command-line tool designed for freelancers and small businesses to transform daily Excel
timesheets into professional, UBL-compliant XML invoices, complete with PDF and HTML/Word renders.

## 🚀 Key Features

- **Automated Calculations**: Aggregates daily hours into project-based invoice lines.
- **UBL 2.1 Compliance**: Generates XML files ready for modern e-invoicing platforms.
- **Multiple Output Formats**:
    - **PDF**: Professional invoice PDF generated from HTML or Word templates.
    - **HTML**: Customizable web-based invoice format.
    - **Word**: Template-based `.docx` generation.
    - **Excel Timesheet**: Fill your existing timesheet template automatically and attach it as a PDF.
- **Smart Configuration**: Hierarchical YAML configuration (Default -> Client -> Project).
- **Interactive UI**: A sleek terminal interface for those who prefer not to type long commands.

## 🛠️ Installation

### Prerequisites

- **Go**: 1.26 or higher.
- **Microsoft Office**: Required for Word/Excel to PDF conversion (Windows only).
- **Google Chrome/Chromium**: Required for HTML to PDF conversion.

### Build from source

```powershell
go build -o ublcli.exe ./cmd/ublcli/main.go
```

## 📖 Quick Start

### 1. Setup Configuration

Create a configuration folder and file in your home directory:
`%USERPROFILE%\.timesheet2ubl\config.yaml`

Example `config.yaml`:

```yaml
default:
  consultant_name: "John Doe"
  daily_rate: 600
  currency: "EUR"
  vat_percent: 21
  supplier:
    name: "My Company SRL"
    company_id: "BE0123456789"
    iban: "BE68..."
    address:
      street: "Rue de la Loi 1"
      city: "Brussels"
      postal_code: "1000"
      country: "Belgium"
```

### 2. Generate an Invoice

Run the following command to process your timesheet:

```powershell
.\ublcli.exe generate --timesheet-in my_march_hours.xlsx --invoice-number 0001
```

### 3. Use the Interactive TUI

Simply run:

```powershell
.\ublcli.exe tui
```

## 📂 Project Structure

- `cmd/ublcli`: CLI Entry point.
- `internal/`: Core logic including Excel parsing, UBL generation, and rendering engines.
- `dist/`: Default output directory for generated invoices.

## 📝 Documentation

For detailed technical information, template placeholders, and advanced configuration, please refer to
the [HANDOFF.md](./HANDOFF.md).

### 🌐 GitHub Pages Deployment

To publish the documentation site (located in `docs/`) to GitHub Pages, follow these steps:

1. **Configure `baseURL`**: In `docs/hugo.toml`, ensure the `baseURL` points to your GitHub Pages URL (e.g.,
   `https://<username>.github.io/<repository>/`).
2. **GitHub Actions**: Create a file at `.github/workflows/hugo.yml` with the following content:

```yaml
name: Deploy Hugo site to Pages

on:
  push:
    branches: [ "main" ]
    paths: [ "docs/**" ]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: false

defaults:
  run:
    working-directory: docs

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          submodules: recursive
          fetch-depth: 0
      - name: Setup Hugo
        uses: peaceiris/actions-hugo@v3
        with:
          hugo-version: 'latest'
          extended: true
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '>=1.20'
          cache-dependency-path: docs/go.mod
      - name: Setup Pages
        id: pages
        uses: actions/configure-pages@v4
      - name: Build with Hugo
        run: |
          hugo mod tidy
          hugo --minify --baseURL "${{ steps.pages.outputs.base_url }}/"
      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: docs/public

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

3. **Enable Pages**: In your GitHub repository settings, go to **Settings > Pages** and select **GitHub Actions** as the
   source for the build and deployment.

## ⚖️ License

[Specify License if applicable, e.g., MIT]
