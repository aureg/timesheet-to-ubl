---
title: CLI
weight: 11
next: usage/config-yaml
---

Cette page documente la commande `generate` de `ublcli`.

## Pré-requis CLI

- **Go 1.26+** pour compiler `ublcli`.
- **Microsoft Office (Excel + Word)** pour les conversions Office -> PDF.
- **Google Chrome** ou **Chromium** pour la conversion HTML -> PDF.

## Build rapide

```powershell
go mod download
go build -o ublcli.exe ./cmd/ublcli/main.go
```

## Commande de base

```powershell
.\ublcli.exe generate --timesheet-in data.xlsx --invoice-number 0001
```

## Référence complète des paramètres (`generate`)

| Option             | Alias       | Type     | Valeurs possibles                                                                    | Défaut                                       | Obligatoire                  |
|--------------------|-------------|----------|--------------------------------------------------------------------------------------|----------------------------------------------|------------------------------|
| `--timesheet-in`   | `--tsin`    | `string` | Chemin vers un fichier Excel lisible par l'outil (`.xlsx`, `.xlsm`...)               | `""`                                         | Oui                          |
| `--sheet`          | -           | `string` | Nom de feuille Excel existant dans le fichier source                                 | `Data`                                       | Non                          |
| `--config`         | -           | `string` | Chemin vers un YAML de configuration                                                 | `%USERPROFILE%\\.timesheet2ubl\\config.yaml` | Oui (avec valeur par défaut) |
| `--template`       | -           | `string` | Chemin explicite (`.html` ou `.docx`) ou nom logique recherché dans `.timesheet2ubl` | `""`                                         | Non                          |
| `--excel-template` | -           | `string` | Chemin vers template Excel (`.xlsx` ou `.xlsm`) ou nom logique                       | `""`                                         | Non                          |
| `--out`            | -           | `string` | Dossier racine de sortie                                                             | `""`                                         | Non                          |
| `--invoice-number` | `--inv-num` | `string` | Exactement 4 caractères (ex: `0001`)                                                 | `""`                                         | Oui                          |

## Règles importantes

- `--timesheet-in`, `--config` et `--invoice-number` sont requis.
- `--invoice-number` est validé sur la longueur (4 caractères).
- Si `--out` est vide, fallback vers `output_dir` de la config, puis `./dist`.
- Le dossier final est toujours suffixé en `INV_<numero>_<client>`.

## Exemple avancé

```powershell
.\ublcli.exe generate `
  --timesheet-in .\input\timesheet-mars.xlsx `
  --sheet Data `
  --config $env:USERPROFILE\.timesheet2ubl\config.yaml `
  --template invoice-nrb `
  --excel-template timesheet-template.xlsx `
  --out .\exports `
  --invoice-number 0001
```

## Voir aussi

- [`Configuration YAML`](/docs/usage-config-yaml)
- [`TUI`](/docs/usage-tui)

