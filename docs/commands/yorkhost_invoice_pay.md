## yorkhost invoice pay

Ouvrir la page de paiement sécurisée d'une facture

### Description

Affiche le statut à jour de la facture et, si elle est impayée, ouvre la page
de paiement de l'espace client (carte avec 3-D Secure, PayPal, crédit du compte).
Aucun moyen de paiement ne transite par le CLI ni par l'API.

```
yorkhost invoice pay <id> [options]
```

### Options

```
  -h, --help   aide de yorkhost invoice pay
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

