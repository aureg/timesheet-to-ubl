---
title: Architecture
weight: 10
---

Le projet `ublcli` est conçu de manière modulaire en Go pour faciliter la maintenance et l'ajout de nouveaux formats de
rendu.

## Structure du code source

L'organisation des packages suit les recommandations standards de Go pour les outils CLI.

- `cmd/ublcli` : Point d'entrée principal de l'application.
- `internal/domain` : Définition des structures de données métier (`Invoice`, `Party`, etc.).
- `internal/excel` : Logique d'importation des fichiers Excel (via `excelize`).
- `internal/invoice` : Calcul des montants, agrégation des lignes et application de la TVA.
- `internal/ubl` : Générateur de XML conforme à la spécification UBL 2.1.
- `internal/render/` : Contient les différents moteurs de rendu :
    - `htmltmpl` : Génération HTML via des templates Go.
    - `word` : Remplissage de placeholders dans des documents `.docx`.
    - `excel` : Remplissage de templates Excel pour validation.
- `internal/pdf` : Wrappers pour transformer les rendus (HTML/Word/Excel) en PDF.
- `internal/orchestrator` : Coordonne tout le flux de travail, de l'import à l'export final.

## Flux de données

1. **Importation** : `internal/excel` lit le fichier Excel et crée une liste brute de prestations.
2. **Traitement** : `internal/invoice` agrège les prestations par projet et calcule les montants selon les tarifs
   configurés.
3. **Génération** : L'orchestrateur appelle simultanément le générateur UBL et les différents moteurs de rendu demandés.
4. **Exportation** : Les fichiers finaux sont sauvegardés et, le cas échéant, les rendus sont convertis en PDF.

## Technologies utilisées

- **Go 1.26+**
- **Excelize** : Pour la manipulation des fichiers Excel.
- **Bubble Tea** : Pour l'interface TUI (Text User Interface).
- **COM Interop (PowerShell)** : Pour l'intégration avec Microsoft Office sous Windows.
