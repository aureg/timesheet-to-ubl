---
title: Configuration YAML
weight: 12
next: usage/tui
---

Cette page détaille le fichier `config.yaml`: clés, valeurs possibles, priorités et exemples.

## Emplacement par défaut

- Windows: `%USERPROFILE%\\.timesheet2ubl\\config.yaml`
- Générique: `~/.timesheet2ubl/config.yaml`

## Structure globale

```yaml
default: { }
clients: { }
projects: { }
client_projects: { }
```

## Section `default`

| Clé                  | Type     | Valeurs possibles                           | Défaut implicite |
|----------------------|----------|---------------------------------------------|------------------|
| `consultant_name`    | `string` | Texte libre                                 | `""`             |
| `currency`           | `string` | Code devise (ex: `EUR`)                     | `""`             |
| `vat_percent`        | `number` | Ex: `0`, `6`, `21`                          | `0`              |
| `daily_rate`         | `number` | Taux journalier                             | `0`              |
| `payment_terms_days` | `int`    | Nombre de jours (ex: `30`)                  | `0`              |
| `invoice_template`   | `string` | Chemin ou nom logique de template facture   | `""`             |
| `excel_template`     | `string` | Chemin ou nom logique de template timesheet | `""`             |
| `output_dir`         | `string` | Dossier racine de sortie                    | `""`             |
| `supplier`           | objet    | Voir tableau ci-dessous                     | objet vide       |

### `default.supplier`

| Clé                   | Type     | Valeurs possibles                | Défaut implicite |
|-----------------------|----------|----------------------------------|------------------|
| `name`                | `string` | Texte libre                      | `""`             |
| `company_id`          | `string` | Numéro d'entreprise / TVA        | `""`             |
| `email`               | `string` | Email                            | `""`             |
| `phone`               | `string` | Téléphone                        | `""`             |
| `iban`                | `string` | IBAN                             | `""`             |
| `bic`                 | `string` | BIC/SWIFT                        | `""`             |
| `address.street`      | `string` | Rue + numéro                     | `""`             |
| `address.city`        | `string` | Ville                            | `""`             |
| `address.postal_code` | `string` | Code postal                      | `""`             |
| `address.country`     | `string` | Code/pays (`BE`, `Belgium`, ...) | `""`             |

## Section `clients`

```yaml
clients:
  "NomClient":
    customer:
      name: "Client SA"
      company_id: "BE0123456789"
      address:
        street: "Rue Exemple 1"
        city: "Bruxelles"
        postal_code: "1000"
        country: "BE"
    vat_percent: 21
    payment_terms_days: 30
    order_reference: "4500001234"
    manager_name: "Jean Manager"
    invoice_template: "invoice-client.html"
    excel_template: "timesheet-client.xlsx"
    output_dir: "./dist-client"
```

| Clé                  | Type     | Effet                                   |
|----------------------|----------|-----------------------------------------|
| `customer`           | objet    | Identité du client facturé              |
| `vat_percent`        | `number` | Surcharge de TVA pour ce client         |
| `payment_terms_days` | `int`    | Surcharge du délai de paiement          |
| `order_reference`    | `string` | Bon de commande injecté dans la facture |
| `manager_name`       | `string` | Manager injecté dans les documents      |
| `invoice_template`   | `string` | Surcharge du template facture           |
| `excel_template`     | `string` | Surcharge du template timesheet         |
| `output_dir`         | `string` | Surcharge du dossier de sortie          |

## Section `projects`

```yaml
projects:
  "PROJ001":
    daily_rate: 750
    invoice_label: "Prestations PROJ001"
```

## Section `client_projects`

- Format de clé: `<NomClient>::<CodeProjet>`

```yaml
client_projects:
  "NomClient::PROJ001":
    daily_rate: 820
    invoice_label: "Prestations spécifiques PROJ001"
```

## Priorité de résolution

1. Base `default`
2. Surcharge `client_projects` (`<client>::<project>`)
3. Surcharge `clients[client]`
4. Surcharge `projects[project]` (si la valeur n'a pas déjà été surchargée)

## Résolution des templates

Ordre de recherche pour la facture:

1. `--template`
2. `invoice_template` de la config résolue
3. Fallback auto `%USERPROFILE%\\.timesheet2ubl\\invoice.html`, puis `%USERPROFILE%\\.timesheet2ubl\\invoice.docx`

Même logique pour `--excel-template` / `excel_template`.

## Exemple complet

```yaml
default:
  consultant_name: "Jean Dupont"
  currency: "EUR"
  vat_percent: 21
  daily_rate: 700
  payment_terms_days: 30
  invoice_template: "invoice.html"
  excel_template: "timesheet.xlsx"
  output_dir: "./dist"
  supplier:
    name: "Ma Societe SRL"
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
  "NRB":
    customer:
      name: "NRB SA"
      company_id: "BE0403201634"
      address:
        street: "Rue des Exemple 10"
        city: "Herstal"
        postal_code: "4040"
        country: "BE"
    vat_percent: 21
    payment_terms_days: 30
    order_reference: "4500012345"
    manager_name: "Manager NRB"
    invoice_template: "invoice-nrb.docx"

projects:
  "CATS":
    daily_rate: 780
    invoice_label: "Prestations CATS"

client_projects:
  "NRB::CATS":
    daily_rate: 820
    invoice_label: "Prestations CATS - NRB"
```

