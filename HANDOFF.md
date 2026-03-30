# Dossier de Passation - ublcli

Ce document contient les instructions et informations nécessaires pour continuer le développement du projet `ublcli` sur un autre ordinateur ou avec une autre IA.

## Présentation du projet
`ublcli` est un outil en Go (v1.26+) permettant de générer des factures conformes au standard UBL XML, ainsi que des rendus HTML et PDF, à partir d'un relevé de prestations journalières au format Excel.

### Fonctionnalités principales :
- **Import Excel** : Lecture flexible (dates FR/EN/Numériques, nombres avec virgules ou points).
- **Logique métier** : Agrégation des prestations par projet, calcul des totaux HT/TVA/TTC.
- **Enrichissement par Config** : Système de résolution de configuration par priorités (Client+Projet > Client > Projet > Défaut) chargé via YAML.
- **Rendu Multi-format** : Génération de HTML via `html/template`, PDF via Chrome Headless, et UBL XML avec PDF embarqué en Base64.
- **Double Interface** : Mode CLI standard et mode TUI interactive (Bubble Tea).

## Architecture technique
Le projet suit une architecture modulaire et testable :

- `cmd/ublcli` : Point d'entrée de l'application.
- `internal/domain` : Structures de données métier (Invoice, Party, etc.).
- `internal/config` : Chargement et résolution de la configuration YAML.
- `internal/excel` : Importateur Excel utilisant `excelize`.
- `internal/invoice` : Service de calcul des montants et agrégation.
- `internal/render/htmltmpl` : Générateur de factures HTML.
- `internal/pdf` : Wrapper pour la conversion HTML -> PDF via Chrome.
- `internal/ubl` : Générateur de XML UBL 2.1.
- `internal/orchestrator` : Coordination du flux complet de génération.
- `internal/cli` : Gestion des commandes et drapeaux (flags).
- `internal/tui` : Interface utilisateur interactive (Bubble Tea).

## Environnement de développement
### Prérequis :
- **Go 1.26** ou supérieur.
- **Google Chrome** ou **Chromium** installé (nécessaire pour la génération PDF).
- Un éditeur compatible Go (JetBrains GoLand recommandé).

### Dépendances majeures :
- `github.com/xuri/excelize/v2` : Lecture Excel.
- `gopkg.in/yaml.v3` : Parsing YAML.
- `github.com/charmbracelet/bubbletea` : Framework TUI.
- `github.com/charmbracelet/lipgloss` : Stylisation TUI.

## Installation et Usage
1. **Initialisation** :
   ```powershell
   go mod download
   ```
2. **Compilation** :
   ```powershell
   go build -o ublcli.exe ./cmd/ublcli/main.go
   ```
3. **Exécution CLI** :
   ```powershell
   .\ublcli.exe generate --excel data.xlsx --config billing.yaml --template invoice.html --out ./dist
   ```
4. **Exécution TUI** :
   ```powershell
   .\ublcli.exe
   ```

## Points de vigilance pour l'IA suivante
1. **Formats de Date/Heure** : L'importateur Excel dans `internal/excel/importer.go` supporte de nombreux formats (`DD/MM/YYYY`, `MM-DD-YY`, etc.) et gère les virgules comme séparateurs décimaux.
2. **Résolution de Configuration** : Si une valeur est manquante dans une section spécifique du YAML, elle remonte vers la valeur par défaut. Voir `internal/config/config.go`.
3. **Gestion du Client** : Le nom du client est prioritairement tiré du fichier de configuration. La colonne "Nom du client" de l'Excel est désormais optionnelle.
4. **PDF via Chrome** : Le chemin vers l'exécutable Chrome est configuré dans `internal/pdf/renderer.go`. Sur Linux, assurez-vous que `google-chrome` ou `chromium` est dans le `PATH`.
5. **Tests** : Les tests unitaires sont cruciaux pour valider la logique de calcul. Exécuter :
   ```powershell
   go test ./internal/...
   ```

## Améliorations futures possibles
- Ajouter le support de multiples devises par ligne de facture (actuellement une seule devise par facture).
- Améliorer le template HTML par défaut pour inclure des logos.
- Ajouter une option pour uploader directement le XML vers une plateforme de facturation (ex: Peppol).
- Gérer les erreurs de conversion PDF de manière plus granulaire si Chrome n'est pas présent.
