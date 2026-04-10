---
title: TUI
weight: 13
next: architecture
---

Cette page couvre le mode interactif de `ublcli` (TUI, Terminal User Interface).

## Lancer la TUI

```powershell
.\ublcli.exe tui
```

## Ce que fait la TUI

- Guide la saisie étape par étape.
- Évite d'écrire une longue commande `generate`.
- Utilise le même moteur de génération que le mode CLI.

## Quand utiliser la TUI

- Vous générez ponctuellement et préférez un assistant interactif.
- Vous ne souhaitez pas mémoriser les options CLI.

## Quand préférer la CLI

- Vous automatisez des exécutions (scripts, CI, tâches récurrentes).
- Vous avez besoin d'une commande reproductible et versionnable.

## Équivalence avec la CLI

La TUI aboutit aux mêmes livrables que `generate`:

- `invoice-ubl.xml`
- `invoice.pdf` (si conversion PDF réussie)
- `invoice.html` ou `invoice_filled.docx`
- `source_<fichier_excel_source>`
- `timesheet_filled.xlsx/.xlsm` et éventuellement `timesheet.pdf`

## Voir aussi

- [`CLI`](/docs/usage/cli)
- [`Configuration YAML`](/docs/usage/config-yaml)

