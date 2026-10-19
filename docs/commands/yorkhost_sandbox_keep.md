## yorkhost sandbox keep

Garder la sandbox : commander un VPS-2 et y copier /home/agent

### Description

Commande un VPS classique (VPS-2, IP publique, géré dans l'espace client).
Une fois la facture payée, le dossier /home/agent de l'essai est copié et vérifié
dans le nouveau VPS, puis l'essai est supprimé. Devis puis confirmation, comme
« yorkhost order create ».

```
yorkhost sandbox keep <id> [options]
```

### Options

```
  -h, --help   aide de yorkhost sandbox keep
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

