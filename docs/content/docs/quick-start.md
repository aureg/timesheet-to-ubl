---
title: Quick Start
weight: 1
next: usage
---

Cette page vous guide pour une première génération en moins de 10 minutes: installation, configuration, templates, puis
exécution de `ublcli`.

## 1) Requirements

- **OS**: Windows recommandé pour les conversions Office -> PDF.
- **Go 1.26+**: nécessaire pour compiler `ublcli`.
- **Microsoft Office (Excel + Word)**: requis pour convertir les fichiers `.xlsx/.xlsm/.docx` en PDF via COM.
- **Google Chrome** ou **Chromium**: requis pour convertir une facture HTML en PDF.

## 2) Installation

Depuis la racine du projet:

```powershell
go mod download
go build -o ublcli.exe ./cmd/ublcli/main.go
```

## 3) Préparer le dossier de configuration

`ublcli` lit par défaut le fichier:

- `%USERPROFILE%\.timesheet2ubl\config.yaml`

Créez le dossier si nécessaire, puis ajoutez au minimum:

```yaml
default:
  consultant_name: "Jean Dupont"
  currency: "EUR"
  vat_percent: 21
  daily_rate: 700
  payment_terms_days: 30

  supplier:
    name: "Ma Société SRL"
    company_id: "BE0123456789"
    email: "facturation@masociete.be"
    phone: "+32 2 000 00 00"
    iban: "BE68000000000000"
    bic: "BBRUBEBB"
    address:
      street: "Rue de l'Exemple 1"
      city: "Bruxelles"
      postal_code: "1000"
      country: "BE"

clients:
  "MonClient":
    customer:
      name: "Mon Client SA"
      company_id: "BE0987654321"
      address:
        street: "Avenue du Client 10"
        city: "Liège"
        postal_code: "4000"
        country: "BE"
    order_reference: "4500012345"
```

## 4) Configurer les templates (HTML, Word, Excel)

### Template facture (`invoice_template` ou `--template`)

Vous pouvez utiliser:

- un chemin explicite (`C:\templates\invoice.docx`, `./invoice.html`),
- ou un nom logique (`invoice-nrb`) résolu dans `%USERPROFILE%\.timesheet2ubl\`.

Configuration type:

```yaml
default:
  invoice_template: "invoice.html"
```

Comportement par défaut si rien n'est fourni:

- l'outil cherche automatiquement `invoice.html`, puis `invoice.docx` dans `%USERPROFILE%\.timesheet2ubl\`.

### Template timesheet Excel (`excel_template` ou `--excel-template`)

Configuration type:

```yaml
default:
  excel_template: "timesheet-template.xlsx"
```

- Valeurs possibles: chemin explicite (`.xlsx`, `.xlsm`) ou nom logique.
- Valeur par défaut: vide (`""`) dans la config.
- Recommandation: renseigner explicitement `excel_template` si vous voulez générer `timesheet_filled.xlsx/.xlsm`.

## 5) Première génération

Exécutez une commande simple:

```powershell
.\ublcli.exe generate --timesheet-in .\input\timesheet.xlsx --invoice-number 0001
```

Exemple plus complet:

```powershell
.\ublcli.exe generate `
  --timesheet-in .\input\timesheet-mars.xlsx `
  --sheet Data `
  --config $env:USERPROFILE\.timesheet2ubl\config.yaml `
  --template invoice `
  --excel-template timesheet-template.xlsx `
  --out .\exports `
  --invoice-number 0001
```

## 6) Vérifier les sorties

Dans le dossier de sortie, `ublcli` crée un sous-dossier:

- `INV_<numero>_<client>`

Vous y trouverez généralement:

- `invoice-ubl.xml`
- `invoice.pdf` (si conversion PDF réussie)
- `invoice.html` ou `invoice_filled.docx`
- `source_<fichier_excel_source>`
- `timesheet_filled.xlsx/.xlsm` et éventuellement `timesheet.pdf` si un template Excel est configuré.

## 7) Aller plus loin

Pour le détail complet, consultez:

- [`Utilisation`](/docs/usage)
- [`CLI`](/docs/usage/cli)
- [`Configuration YAML`](/docs/usage/config-yaml)
- [`TUI`](/docs/usage/tui)
