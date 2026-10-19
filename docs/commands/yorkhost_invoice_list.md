## yorkhost invoice list

Lister les factures

```
yorkhost invoice list [options]
```

### Exemples

```
  yorkhost invoice list --status unpaid
```

### Options

```
  -h, --help            aide de yorkhost invoice list
      --limit int       nombre maximum de résultats
      --status string   filtrer : paid, unpaid, cancelled, refunded
```

### Options globales

```
      --api-url string   URL de l'API (env YORKHOST_API_URL)
      --debug            journal des requêtes HTTP sur stderr (sans secrets)
      --no-color         désactiver les couleurs (env NO_COLOR)
  -o, --output string    format de sortie : table, wide, json, yaml (env YORKHOST_OUTPUT)
  -p, --profile string   profil de configuration (env YORKHOST_PROFILE)
  -q, --quiet            n'afficher que le résultat
  -y, --yes              répondre oui aux confirmations (env YORKHOST_YES=1)
```

### Voir aussi

* [yorkhost invoice](yorkhost_invoice.md)	 - Factures : liste, détail, PDF, paiement

