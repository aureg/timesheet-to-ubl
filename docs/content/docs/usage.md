---
title: Utilisation
weight: 10
---

# Utilisation du logiciel

Cette page décrit comment installer et utiliser `ublcli` pour générer vos factures.

## Pré-requis

- **Go 1.26** ou supérieur.
- **Microsoft Office** (Excel/Word) : Nécessaire pour les conversions vers PDF via COM (Windows).
- **Google Chrome** ou **Chromium** : Pour les conversions HTML vers PDF.

## Installation

```powershell
go mod download
go build -o ublcli.exe ./cmd/ublcli/main.go
```

## Commandes de base

### Mode CLI

```powershell
.\ublcli.exe generate --timesheet-in data.xlsx --invoice-number 2026-0001
```

Options principales :

- `--tsin` ou `--timesheet-in` : Le fichier Excel source.
- `--inv-num` ou `--invoice-number` : Numéro de la facture.
- `--template` : Chemin vers un template Word ou HTML spécifique.

### Mode Interactif (TUI)

```powershell
.\ublcli.exe tui
```

L'interface graphique dans le terminal vous guidera étape par étape.

## Configuration

La configuration par défaut se trouve dans `%userHome%\.timesheet2ubl\config.yaml`.
Elle permet de définir les clients, les projets et les tarifs journaliers par défaut.
