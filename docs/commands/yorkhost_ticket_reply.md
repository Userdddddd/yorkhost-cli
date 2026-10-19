## yorkhost ticket reply

Répondre à un ticket

```
yorkhost ticket reply <numéro|id> [options]
```

### Exemples

```
  yorkhost ticket reply 10234 --body "C'est réglé, merci !"
```

### Options

```
  -b, --body string   message (texte, @fichier ou - pour stdin)
  -h, --help          aide de yorkhost ticket reply
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

* [yorkhost ticket](yorkhost_ticket.md)	 - Support : tickets, réponses

