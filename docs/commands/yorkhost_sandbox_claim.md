## yorkhost sandbox claim

Rattacher un essai anonyme (ssh yorkhost.sh) et le passer en VPS

### Description

Rattache à votre compte un essai créé anonymement avec « ssh yorkhost.sh », à partir
du jeton yhsc_… ou du lien « Keep it » affiché dans l'essai, et commande le VPS-2 qui
le recevra (devis puis confirmation).

```
yorkhost sandbox claim <jeton|lien> [options]
```

### Options

```
  -h, --help   aide de yorkhost sandbox claim
      --pay    ouvrir la page de paiement si la facture est impayée
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

* [yorkhost sandbox](yorkhost_sandbox.md)	 - Sandboxes agents : VM d'essai jetables (SSH via yorkhost.sh)

