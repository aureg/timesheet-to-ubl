# Dossier de Passation - ublcli

Ce document contient les instructions et informations nécessaires pour continuer le développement du projet `ublcli` sur un autre ordinateur ou avec une autre IA.

## Présentation du projet
`ublcli` est un outil en Go (v1.26+) permettant de générer des factures conformes au standard UBL XML, ainsi que des rendus HTML et PDF, à partir d'un relevé de prestations journalières au format Excel.

### Fonctionnalités principales :
- **Import Excel** : Lecture flexible (dates FR/EN/Numériques, nombres avec virgules ou points).
- **Logique métier** : Agrégation des prestations par projet, calcul des totaux HT/TVA/TTC.
- **Enrichissement par Config** : Système de résolution de configuration par priorités (Client+Projet > Client > Projet > Défaut) chargé via YAML.
- **Rendu Multi-format** : Génération de HTML (via `html/template`) ou Word (via placeholders `{{Key}}`), conversion PDF
  via Chrome (HTML) ou LibreOffice (Excel/Word), et UBL XML.
- **Double Interface** : Mode CLI standard et mode TUI interactive (Bubble Tea).

## Architecture technique
Le projet suit une architecture modulaire et testable :

- `cmd/ublcli` : Point d'entrée de l'application.
- `internal/domain` : Structures de données métier (Invoice, Party, etc.).
- `internal/config` : Chargement et résolution de la configuration YAML.
- `internal/excel` : Importateur Excel utilisant `excelize`.
- `internal/invoice` : Service de calcul des montants et agrégation.
- `internal/render/htmltmpl` : Générateur de factures HTML via `html/template`.
- `internal/render/word` : Générateur de factures via templates Word (.docx).
- `internal/pdf` : Wrapper pour la conversion HTML -> PDF via Chrome.
- `internal/ubl` : Générateur de XML UBL 2.1.
- `internal/orchestrator` : Coordination du flux complet de génération.
- `internal/cli` : Gestion des commandes et drapeaux (flags).
- `internal/tui` : Interface utilisateur interactive (Bubble Tea).

## Environnement de développement
### Prérequis :
- **Go 1.26** ou supérieur.
- **LibreOffice** : Installé dans `C:\Program Files\LibreOffice` (nécessaire pour la conversion Excel/Word -> PDF).
- **Google Chrome** ou **Chromium** : Installé (pour la conversion HTML -> PDF).
- Un éditeur compatible Go (JetBrains GoLand recommandé).

### Dépendances majeures :

- `github.com/nguyenthenguyen/docx` : Manipulation de fichiers Word.
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
   # Avec template HTML
   .\ublcli.exe generate --excel data.xlsx --config billing.yaml --template invoice.html --out ./dist
   # Avec template Word
   .\ublcli.exe generate --excel data.xlsx --config billing.yaml --template invoice.docx --out ./dist
   ```
4. **Exécution TUI** :
   ```powershell
   .\ublcli.exe
   ```

## Documentation des Templates

L'outil supporte deux types de templates : HTML (`.html`) et Word (`.docx`). Les données de la facture sont injectées
différemment selon le format.

### Placeholders Communs (HTML & Word)

Ces valeurs sont disponibles sous forme de chaînes de caractères formatées :

| Placeholder          | Description                       | Format / Exemple   |
|:---------------------|:----------------------------------|:-------------------|
| `{{InvoiceNumber}}`  | Numéro de la facture              | `2026-001`         |
| `{{InvoiceDate}}`    | Date d'émission                   | `31/03/2026`       |
| `{{DueDate}}`        | Date d'échéance                   | `30/04/2026`       |
| `{{PeriodStart}}`    | Début de la période de prestation | `01/03/2026`       |
| `{{PeriodEnd}}`      | Fin de la période de prestation   | `31/03/2026`       |
| `{{ClientName}}`     | Nom du client (config)            | `Nom du Client`    |
| `{{TotalHours}}`     | Somme totale des heures           | `145.50`           |
| `{{Subtotal}}`       | Montant total HT                  | `9457.50`          |
| `{{VATAmount}}`      | Montant total de la TVA           | `1986.08`          |
| `{{TotalAmount}}`    | Montant total TTC                 | `11443.58`         |
| `{{Currency}}`       | Devise utilisée                   | `EUR`              |
| `{{OrderReference}}` | Numéro de bon de commande (PO)    | `PO-12345`         |
| `{{Now}}`            | Date et heure de génération       | `31/03/2026 11:45` |

### Spécificités du Template HTML (Go `html/template`)

En plus des placeholders ci-dessus, le moteur Go expose les structures de données complètes du domaine. Vous pouvez
utiliser des boucles et des accès aux propriétés imbriquées.

**Objets disponibles :**

- `{{.Supplier}}` : Objet `Party` du fournisseur (Name, CompanyID, Email, Phone, IBAN, BIC, Address).
- `{{.Customer}}` : Objet `Party` du client (Name, CompanyID, Email, Phone, Address).
- `{{.Lines}}` : Liste des lignes de facture (`InvoiceLine`).

**Exemple de boucle sur les lignes :**

```html
{{range .Lines}}
<tr>
   <td>{{.Description}}</td>
   <td>{{printf "%.2f" .Quantity}}</td>
   <td>{{printf "%.2f" .UnitPrice}}</td>
   <td>{{printf "%.2f" .NetAmount}}</td>
</tr>
{{end}}
```

### Spécificités du Template Word (`.docx`)

Le support Word utilise un remplacement de texte simple. Les données complexes sont pré-formatées en chaînes de
caractères.

**Placeholders additionnels :**

- `{{SupplierName}}` : Nom du fournisseur.
- `{{CustomerName}}` : Nom du client.
- `{{Lines}}` : Un résumé textuel multi-lignes de toutes les prestations (Description, Quantité, Prix unitaire, Total).

*Note : Pour le format Word, les boucles dans les tableaux ne sont pas encore supportées. Le placeholder `{{Lines}}`
insère un bloc de texte brut récapitulatif.*

## Points de vigilance pour l'IA suivante
1. **Formats de Date/Heure** : L'importateur Excel dans `internal/excel/importer.go` supporte de nombreux formats (`DD/MM/YYYY`, `MM-DD-YY`, etc.) et gère les virgules comme séparateurs décimaux.
2. **Résolution de Configuration** : Si une valeur est manquante dans une section spécifique du YAML, elle remonte vers la valeur par défaut. Voir `internal/config/config.go`.
3. **Gestion du Client** : Le nom du client est prioritairement tiré du fichier de configuration. La colonne "Nom du client" de l'Excel est désormais optionnelle.
4. **PDF via Chrome & LibreOffice** :
   - La conversion HTML -> PDF utilise Chrome.
   - La conversion Excel/Word -> PDF utilise LibreOffice (`soffice.exe`).
   - Assurez-vous que les chemins dans `internal/pdf/renderer.go` correspondent à votre installation.
5. **Tests** : Les tests unitaires sont cruciaux pour valider la logique de calcul. Exécuter :
   ```powershell
   go test ./internal/...
   ```

## Améliorations futures possibles
- Ajouter le support de multiples devises par ligne de facture (actuellement une seule devise par facture).
- Améliorer le template HTML par défaut pour inclure des logos.
- Ajouter une option pour uploader directement le XML vers une plateforme de facturation (ex: Peppol).
- Gérer les erreurs de conversion PDF de manière plus granulaire si Chrome n'est pas présent.
